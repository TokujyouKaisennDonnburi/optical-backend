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
func (r *TokenRedisRepository) AddToWhitelist(ctx context.Context, refreshToken *user.RefreshToken) error {
	key := r.getWhitelistKey(refreshToken.UserId)

	// 既存のリストを取得
	entries, err := r.getWhitelistEntries(ctx, key)
	if err != nil {
		return err
	}

	// 期限切れトークンを除外
	entries = filterValidEntries(entries)

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

	// 180日のTTLを設定
	ttl := time.Second * time.Duration(user.REFRESH_TOKEN_EXPIRE)
	return r.client.Set(ctx, key, data, ttl).Err()
}

// 指定されたトークンがホワイトリストに存在するか確認する
func (r *TokenRedisRepository) IsWhitelisted(ctx context.Context, userId uuid.UUID, tokenId uuid.UUID) error {
	key := r.getWhitelistKey(userId)

	entries, err := r.getWhitelistEntries(ctx, key)
	if err != nil {
		return err
	}

	// 期限切れトークンを除外
	validEntries := filterValidEntries(entries)

	// 期限切れトークンがあれば保存し直す
	if len(validEntries) != len(entries) {
		data, err := json.Marshal(validEntries)
		if err != nil {
			return err
		}
		ttl := time.Second * time.Duration(user.REFRESH_TOKEN_EXPIRE)
		if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
			return err
		}
	}

	for _, entry := range validEntries {
		if entry.TokenId == tokenId.String() {
			// 成功時もTTLを延長
			ttl := time.Second * time.Duration(user.REFRESH_TOKEN_EXPIRE)
			_ = r.client.Expire(ctx, key, ttl)
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

// 期限切れでない有効なエントリのみを返す
func filterValidEntries(entries []tokenWhitelistEntry) []tokenWhitelistEntry {
	now := time.Now().UTC()
	valid := make([]tokenWhitelistEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.ExpiresIn.After(now) {
			valid = append(valid, entry)
		}
	}
	return valid
}

