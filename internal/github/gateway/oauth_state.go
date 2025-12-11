package gateway

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (r *GithubApiRepository) SaveOauthState(
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
