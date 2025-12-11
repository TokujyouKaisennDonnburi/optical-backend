package gateway

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (r *GithubApiRepository) LinkUser(
	ctx context.Context,
	code, state string,
) error {
	result, err := r.redisClient.GetDel(ctx, getOauthStateKey(state)).Result()
	if err != nil {
		return err
	}
	userId, err := uuid.Parse(result)
	if err != nil {
		return err
	}
	accessToken, err := postOauthAccessToken(code)
	if err != nil {
		return err
	}
	githubResp, err := postGithubUser(accessToken)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO user_githubs(user_id, github_id, github_name, created_at, updated_at)
			VALUES(:userId, :githubId, :githubName, :createdAt, :updatedAt)
		ON CONFLICT(user_id) 
		DO UPDATE SET
			github_id = :githubId,
			github_name = :githubName,
			updated_at = :updatedAt
	`
	_, err = r.db.NamedExecContext(ctx, query, map[string]any{
		"userId":     userId,
		"githubId":   githubResp.Id,
		"githubName": githubResp.Login,
		"createdAt":  time.Now(),
		"updatedAt":  time.Now(),
	})
	return err
}
