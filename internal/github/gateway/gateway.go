package gateway

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type GithubApiRepository struct {
	db                          *sqlx.DB
	installationIdEncryptionKey []byte
}

type StateRedisRepository struct {
	db                *sqlx.DB
	redisClient       *redis.Client
	redisEncryptionKey []byte
}

func NewGithubApiRepository(db *sqlx.DB, installationIdEncryptionKey []byte) *GithubApiRepository {
	if db == nil {
		panic("db is nil")
	}
	if installationIdEncryptionKey == nil {
		panic("installationIdEncryptionKey is nil")
	}
	return &GithubApiRepository{
		db:                          db,
		installationIdEncryptionKey: installationIdEncryptionKey,
	}
}

func NewStateRedisRepository(
	db *sqlx.DB,
	redisClient *redis.Client,
	redisEncryptionKey []byte,
) *StateRedisRepository {
	if db == nil {
		panic("db is nil")
	}
	if redisClient == nil {
		panic("redisClient is nil")
	}
	if redisEncryptionKey == nil {
		panic("redisEncryptionKey is nil")
	}
	return &StateRedisRepository{
		db:                 db,
		redisClient:        redisClient,
		redisEncryptionKey: redisEncryptionKey,
	}
}
