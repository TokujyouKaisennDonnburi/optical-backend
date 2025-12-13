package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
	"github.com/google/uuid"
)

func (r *GithubApiRepository) LinkUser(
	ctx context.Context,
	userId uuid.UUID,
	code string,
) error {
	accessToken, err := api.PostOauthAccessToken(code)
	if err != nil {
		return err
	}
	githubResp, err := api.GetGithubUser(accessToken)
	if err != nil {
		return err
	}
	githubEmail, err := api.GetGithubPrimaryEmail(accessToken)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO user_githubs(user_id, github_id, github_name, github_email, sso_login, created_at, updated_at)
			VALUES(:userId, :githubId, :githubName, :githubEmail, sso_login, :createdAt, :updatedAt)
		ON CONFLICT(user_id) 
		DO UPDATE SET
			github_id = :githubId,
			github_name = :githubName,
			github_email = :githubEmail,
			updated_at = :updatedAt
	`
	_, err = r.db.NamedExecContext(ctx, query, map[string]any{
		"userId":      userId,
		"githubId":    githubResp.Id,
		"githubName":  githubResp.Login,
		"githubEmail": githubEmail,
		"createdAt":   time.Now(),
		"updatedAt":   time.Now(),
	})
	return err
}
