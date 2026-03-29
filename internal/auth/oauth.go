package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	authorizeURL = "https://www.strava.com/oauth/authorize"
	tokenURL     = "https://www.strava.com/oauth/token"
	redirectURI  = "http://localhost:8089/callback"
)

func Authenticate(clientID, clientSecret string) (*Token, error) {
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	listener, err := net.Listen("tcp", "localhost:8089")
	if err != nil {
		return nil, fmt.Errorf("could not start local server: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no authorization code received")
			fmt.Fprintln(w, "Error: no authorization code received. Please try again.")
			return
		}
		codeCh <- code
		fmt.Fprintln(w, "Authentication successful! You can close this tab and return to your terminal.")
	})

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=activity:read_all",
		authorizeURL, clientID, url.QueryEscape(redirectURI))

	fmt.Println("Open this URL in your browser to authorize z2-cli:")
	fmt.Println()
	fmt.Println(authURL)
	fmt.Println()
	fmt.Println("Waiting for authorization...")

	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		server.Shutdown(context.Background())
		return nil, err
	case <-time.After(2 * time.Minute):
		server.Shutdown(context.Background())
		return nil, fmt.Errorf("authorization timed out after 2 minutes")
	}

	server.Shutdown(context.Background())

	token, err := exchangeCode(clientID, clientSecret, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func exchangeCode(clientID, clientSecret, code string) (*Token, error) {
	data := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange returned status %d", resp.StatusCode)
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresAt    int64  `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("could not parse token response: %w", err)
	}

	return &Token{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    result.ExpiresAt,
	}, nil
}

func RefreshAccessToken(clientID, clientSecret string, token *Token) (*Token, error) {
	data := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"refresh_token": {token.RefreshToken},
		"grant_type":    {"refresh_token"},
	}

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh returned status %d", resp.StatusCode)
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresAt    int64  `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("could not parse refresh response: %w", err)
	}

	return &Token{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    result.ExpiresAt,
	}, nil
}
