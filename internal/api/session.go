package api

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

var sessionKey []byte

func init() {
	sessionKey = make([]byte, 32)
	if _, err := rand.Read(sessionKey); err != nil {
		panic("failed to generate session key: " + err.Error())
	}
}

const sessionCookieName = "z2_session"

// signSession generates a session token as "nonce.signature" where the nonce
// is random and the signature is HMAC-SHA256(nonce). Each call produces a
// unique token.
func signSession() string {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		panic("failed to generate session nonce: " + err.Error())
	}
	nonceHex := hex.EncodeToString(nonce)

	mac := hmac.New(sha256.New, sessionKey)
	mac.Write([]byte(nonceHex))
	sig := hex.EncodeToString(mac.Sum(nil))

	return nonceHex + "." + sig
}

func setSessionCookie(w http.ResponseWriter, r *http.Request) {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    signSession(),
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   scheme == "https",
		SameSite: http.SameSiteLaxMode,
	})
}

func validSession(r *http.Request) bool {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return false
	}
	parts := strings.SplitN(cookie.Value, ".", 2)
	if len(parts) != 2 {
		return false
	}
	mac := hmac.New(sha256.New, sessionKey)
	mac.Write([]byte(parts[0]))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(parts[1]), []byte(expected))
}

func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !validSession(r) {
			writeError(w, http.StatusUnauthorized, "not authenticated")
			return
		}
		next.ServeHTTP(w, r)
	})
}
