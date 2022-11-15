package repository

import (
	"context"
	models "github.com/Vitaly-Baidin/auth-api/internal/model/db"
)

type TokenRepository interface {
	Store(ctx context.Context, token *models.Token) error
	Drop(ctx context.Context, userID uint, tokenType string) error
	GetByKey(ctx context.Context, userID uint, tokenType string) (*models.Token, error)
}
