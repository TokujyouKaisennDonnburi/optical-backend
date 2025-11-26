package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	REDIS_TOKEN_WHITELIST_NAME = "whitelist:token:list"
)

type TokenRedisRepository struct {
	client *redis.Client
}

func NewTokenRedisRepository(client *redis.Client) *TokenRedisRepository {
	if client == nil {
		panic("Redis client is nil")
	}
	return &TokenRedisRepository{
		client: client,
	}
}

func (r *TokenRedisRepository) AddToWhitelist(refreshToken *user.RefreshToken) error {
	result := r.client.SAdd(context.Background(), REDIS_TOKEN_WHITELIST_NAME, refreshToken.Id.String())
	return result.Err()
}

func (r *TokenRedisRepository) IsWhitelisted(tokenId uuid.UUID) error {
	exists, err := r.client.SIsMember(context.Background(), REDIS_TOKEN_WHITELIST_NAME, tokenId.String()).Result()
	if err != nil {
		return err
	}
	if !exists {
		return apperr.ValidationError(tokenId.String() + " is not in whitelist")
	}
	return nil
}
