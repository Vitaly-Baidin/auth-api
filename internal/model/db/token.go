package models

import (
	"time"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Token struct {
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (model *Token) GetResponseJson() map[string]any {
	m := map[string]any{
		"token":   model.Token,
		"expires": model.ExpiresAt.Format("2006-01-02 15:04:05"),
	}
	return m
}

func NewToken(userID uint, tokenString string, tokenType string, expiresAt time.Time) *Token {
	return &Token{
		UserID:    userID,
		Token:     tokenString,
		Type:      tokenType,
		ExpiresAt: expiresAt,
	}
}
