package gateway

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type GithubApiRepository struct {
	db *sqlx.DB
}

type StateRedisRepository struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewGithubApiRepository(db *sqlx.DB) *GithubApiRepository {
	if db == nil {
		panic("db is nil")
	}
	return &GithubApiRepository{
		db: db,
	}
}

func NewStateRedisRepository(db *sqlx.DB, redisClient *redis.Client) *StateRedisRepository {
	if db == nil {
		panic("db is nil")
	}
	if redisClient == nil {
		panic("redisClient is nil")
	}
	return &StateRedisRepository{
		db:          db,
		redisClient: redisClient,
	}
}
