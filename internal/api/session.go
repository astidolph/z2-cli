package api

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

var sessionKey []byte

func init() {
	sessionKey = make([]byte, 32)
	if _, err := rand.Read(sessionKey); err != nil {
		panic("failed to generate session key: " + err.Error())
	}
}

const sessionCookieName = "z2_session"

func signSession() string {
	mac := hmac.New(sha256.New, sessionKey)
	mac.Write([]byte("z2-authenticated"))
	return hex.EncodeToString(mac.Sum(nil))
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
	return hmac.Equal([]byte(cookie.Value), []byte(signSession()))
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
