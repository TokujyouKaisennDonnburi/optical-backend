package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (r *GithubApiRepository) InstallToCalendar(
	ctx context.Context,
	userId, calendarId uuid.UUID,
	code, installationId string,
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		response, err := api.GetGithubInstallation(ctx, installationId)
		if err != nil {
			return err
		}
		query := `
			INSERT INTO calendar_githubs(calendar_id, github_id, github_name, installation_id, created_at, updated_at)
				VALUES(:calendarId, :githubId, :githubName, :installationId, :createdAt, :updatedAt)
			ON CONFLICT(calendar_id) 
			DO UPDATE SET
				github_id = :githubId,
				github_name = :githubName,
				installation_id = :installationId,
				updated_at = :updatedAt
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"calendarId":     calendarId,
			"githubId":       response.Account.Id,
			"githubName":     response.Account.Login,
			"installationId": installationId,
			"createdAt":      time.Now(),
			"updatedAt":      time.Now(),
		})
		if err != nil {
			return err
		}
		accessToken, err := api.PostOauthAccessToken(code)
		if err != nil {
			return err
		}
		githubUser, err := api.GetGithubUser(accessToken)
		if err != nil {
			return err
		}
		query = `
			INSERT INTO user_githubs(user_id, github_id, github_name, github_email, sso_login, created_at, updated_at)
				VALUES(:userId, :githubId, :githubName, :githubEmail, false, :createdAt, :updatedAt)
			ON CONFLICT(user_id) 
			DO UPDATE SET
				github_id = :githubId,
				github_name = :githubName,
				github_email = :githubEmail,
				updated_at = :updatedAt
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"userId":      userId,
			"githubId":    githubUser.Id,
			"githubName":  githubUser.Name,
			"githubEmail": githubUser.Email,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		})
		return err
	})
}
