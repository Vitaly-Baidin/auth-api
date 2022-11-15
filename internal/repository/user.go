package repository

import (
	"context"
	"errors"
	models "github.com/Vitaly-Baidin/auth-api/internal/model/db"
)

type UserRepository interface {
	Store(ctx context.Context, user *models.User) (uint, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	IfExistsByEmail(ctx context.Context, email string) error
	IfExistsByLogin(ctx context.Context, login string) error
}

var (
	ErrUserExists  = errors.New("user already exists")
	ErrEmailExists = errors.New("email already exists")
)
