package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	REDIS_TOKEN_WHITELIST_PREFIX = "token:whitelist:"
)

//  ホワイトリストに保存するトークン情報を表す
type tokenWhitelistEntry struct {
	TokenId   string    `json:"tokenId"`
	ExpiresIn time.Time `json:"expiresIn"`
}

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

// ユーザーIDからRedisのキーを生成する
func (r *TokenRedisRepository) getWhitelistKey(userId uuid.UUID) string {
	return fmt.Sprintf("%s%s", REDIS_TOKEN_WHITELIST_PREFIX, userId.String())
}

// リフレッシュトークンをホワイトリストに追加する
func (r *TokenRedisRepository) AddToWhitelist(refreshToken *user.RefreshToken) error {
	ctx := context.Background()
	key := r.getWhitelistKey(refreshToken.UserId)

	// 既存のリストを取得
	entries, err := r.getWhitelistEntries(ctx, key)
	if err != nil {
		return err
	}

	// 新しいエントリを追加
	newEntry := tokenWhitelistEntry{
		TokenId:   refreshToken.Id.String(),
		ExpiresIn: refreshToken.ExpiresIn,
	}
	entries = append(entries, newEntry)

	// JSONにシリアライズして保存
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, 0).Err()
}

// 指定されたトークンがホワイトリストに存在するか確認する
func (r *TokenRedisRepository) IsWhitelisted(userId uuid.UUID, tokenId uuid.UUID) error {
	ctx := context.Background()
	key := r.getWhitelistKey(userId)

	entries, err := r.getWhitelistEntries(ctx, key)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.TokenId == tokenId.String() {
			return nil
		}
	}

	return apperr.UnauthorizedError(tokenId.String() + " is not in whitelist")
}

// Redisからホワイトリストのエントリ一覧を取得する
func (r *TokenRedisRepository) getWhitelistEntries(ctx context.Context, key string) ([]tokenWhitelistEntry, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return []tokenWhitelistEntry{}, nil
	}
	if err != nil {
		return nil, err
	}

	var entries []tokenWhitelistEntry
	if err := json.Unmarshal([]byte(data), &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

