package gateway

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

const (
	GITHUB_BASE_URL = "https://api.github.com"
)

type GithubApiRepository struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewGithubApiRepository(db *sqlx.DB, redisClient *redis.Client) *GithubApiRepository {
	if db == nil {
		panic("db is nil")
	}
	if redisClient == nil {
		panic("redisClient is nil")
	}
	return &GithubApiRepository{
		db:          db,
		redisClient: redisClient,
	}
}
