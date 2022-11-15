package repository

import (
	"context"
	"fmt"
	models "github.com/Vitaly-Baidin/auth-api/internal/model/db"
	"github.com/Vitaly-Baidin/auth-api/pkg/redis"
	"github.com/go-redis/cache/v8"
	"time"
)

type TokenRepoRedis struct {
	*redis.Redis
}

func NewTokenRepoRedis(r *redis.Redis) *TokenRepoRedis {
	return &TokenRepoRedis{Redis: r}
}

func (r *TokenRepoRedis) Store(ctx context.Context, token *models.Token) error {
	key := fmt.Sprintf("%d:%s", token.UserID, token.Type)

	err := r.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: token.Token,
		TTL:   time.Since(token.ExpiresAt).Abs(),
	})
	if err != nil {
		return fmt.Errorf("failed store token: %v", err)
	}

	return nil
}

func (r *TokenRepoRedis) Drop(ctx context.Context, userID uint, tokenType string) error {
	key := fmt.Sprintf("%d:%s", userID, tokenType)

	err := r.Cache.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed drop token: %v", err)
	}

	return nil
}

func (r *TokenRepoRedis) GetByKey(ctx context.Context, userID uint, tokenType string) (*models.Token, error) {
	var t models.Token
	key := fmt.Sprintf("%d:%s", userID, tokenType)

	t.UserID = userID
	t.Type = tokenType

	err := r.Cache.Get(ctx, key, &t.Token)
	if err != nil {
		return nil, fmt.Errorf("failed get token by key: %v", err)
	}

	return &t, nil
}
