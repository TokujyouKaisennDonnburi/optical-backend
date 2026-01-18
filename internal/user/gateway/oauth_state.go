package gateway

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	GOOGLE_OAUTH_STATE_KEY_PREFIX = "google:oauth:state:"
)

type GoogleOauthStateRedisRepository struct {
	clientId     string
	clientSecret string
	redirectUri  string
	redisClient  *redis.Client
}

func NewGoogleOauthStateRedisRepository(clientId, clientSecret, redirectUri string, redisClient *redis.Client) *GoogleOauthStateRedisRepository {
	return &GoogleOauthStateRedisRepository{
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectUri:  redirectUri,
		redisClient:  redisClient,
	}
}

func getGoogleOauthStateKey(state string) string {
	return GOOGLE_OAUTH_STATE_KEY_PREFIX + state
}

func (r *GoogleOauthStateRedisRepository) SaveOauthState(
	ctx context.Context,
	userId uuid.UUID,
	state string,
) error {
	exp := time.Duration(time.Minute * 10)
	status := r.redisClient.SetEx(ctx, getGoogleOauthStateKey(state), userId.String(), exp)
	return status.Err()
}

func (r *GoogleOauthStateRedisRepository) GetOauthState(
	ctx context.Context,
	state string,
) (uuid.UUID, error) {
	result, err := r.redisClient.GetDel(ctx, getGoogleOauthStateKey(state)).Result()
	if err != nil {
		return uuid.Nil, err
	}
	userId, err := uuid.Parse(result)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}

func (r *GoogleOauthStateRedisRepository) GetClientId() string {
	return r.clientId
}

func (r *GoogleOauthStateRedisRepository) GetRedirectUri() string {
	return r.redirectUri
}
