package auth

import (
	"github.com/z2-cli/internal/model"
	"github.com/z2-cli/internal/storage"
)

// Token is an alias for model.Token, preserving backward compatibility.
type Token = model.Token

func SaveToken(token *Token) error {
	return storage.Get().SaveToken(token)
}

func LoadToken() (*Token, error) {
	return storage.Get().LoadToken()
}
