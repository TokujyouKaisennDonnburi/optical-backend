package gateway

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (r *StateRedisRepository) GetOauthState(
	ctx context.Context,
	state string,
) (uuid.UUID, error) {
	result, err := r.redisClient.GetDel(ctx, getOauthStateKey(state)).Result()
	if err != nil {
		return uuid.Nil, err
	}
	userId, err := uuid.Parse(result)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}

func (r *StateRedisRepository) SaveOauthState(
	ctx context.Context,
	userId uuid.UUID,
	state string,
) error {
	exp := time.Duration(time.Minute * 10)
	status := r.redisClient.SetEx(ctx, getOauthStateKey(state), userId.String(), exp)
	return status.Err()
}

func getOauthStateKey(state string) string {
	return "github:oauth:state:" + state
}
