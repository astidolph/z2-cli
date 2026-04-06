package model

import "time"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func (t *Token) IsExpired() bool {
	return time.Now().Unix() >= t.ExpiresAt
}
