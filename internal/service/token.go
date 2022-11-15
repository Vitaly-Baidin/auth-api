package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Vitaly-Baidin/auth-api/config"
	db "github.com/Vitaly-Baidin/auth-api/internal/model/db"
	"github.com/Vitaly-Baidin/auth-api/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

var (
	ErrValidToken   = errors.New("not valid token")
	ErrExpiredToken = errors.New("token is expired")
)

type TokenService interface {
	CreateToken(ctx context.Context, user *db.User, tokenType string, expiresAt time.Time) (*db.Token, error)
	DeleteToken(ctx context.Context, userID uint, tokenType string) error
	GenerateAccessTokens(ctx context.Context, user *db.User) (*db.Token, *db.Token, error)
	VerifyToken(ctx context.Context, token string, tokenType string) (*db.Token, error)
}

type TokenServ struct {
	cfg  *config.Config
	repo repository.TokenRepository
}

func NewTokenService(cfg *config.Config, r repository.TokenRepository) *TokenServ {
	return &TokenServ{cfg: cfg, repo: r}
}

func (s *TokenServ) CreateToken(ctx context.Context, u *db.User, tokenType string, expiresAt time.Time) (*db.Token, error) {
	claims := &db.UserClaims{
		Email: u.Email,
		Type:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   strconv.Itoa(int(u.ID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed create access token: %v", err)
	}

	tokenModel := db.NewToken(u.ID, tokenString, tokenType, expiresAt)
	err = s.repo.Store(ctx, tokenModel)
	if err != nil {
		return nil, fmt.Errorf("failed save access token to db: %v", err)
	}

	return tokenModel, nil
}

func (s *TokenServ) DeleteToken(ctx context.Context, userID uint, tokenType string) error {
	err := s.repo.Drop(ctx, userID, tokenType)
	if err != nil {
		return fmt.Errorf("failed delete token: %v", err)
	}

	return nil
}

func (s *TokenServ) GenerateAccessTokens(ctx context.Context, u *db.User) (*db.Token, *db.Token, error) {
	accessExpiresAt := time.Now().Add(time.Duration(s.cfg.JWT.AccessExpireInMinute) * time.Minute)
	refreshExpiresAt := time.Now().Add(time.Duration(s.cfg.JWT.RefreshExpireInHour) * time.Hour * 24)

	accessToken, err := s.CreateToken(ctx, u, db.TokenTypeAccess, accessExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.CreateToken(ctx, u, db.TokenTypeRefresh, refreshExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *TokenServ) VerifyToken(ctx context.Context, token string, tokenType string) (*db.Token, error) {
	claims := &db.UserClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil || claims.Type != tokenType {
		return nil, ErrValidToken
	}

	if time.Now().Sub(claims.ExpiresAt.Time) > 10*time.Second {
		return nil, ErrExpiredToken
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed convert subjects to int: %v", err)
	}

	key, err := s.repo.GetByKey(ctx, uint(userID), tokenType)
	if err != nil {
		return nil, fmt.Errorf("failed get key: %v", err)
	}

	return key, nil
}
