package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InstallationGetResponse struct {
	Id      int                            `json:"id"`
	Account InstallationGetResponseAccount `json:"account"`
}

type InstallationGetResponseAccount struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

func (r *GithubApiRepository) InstallToCalendar(
	ctx context.Context,
	userId, calendarId uuid.UUID,
	code, installationId string,
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		client := http.Client{}
		req, err := http.NewRequestWithContext(
			ctx,
			"GET",
			GITHUB_BASE_URL+"/app/installations/"+installationId,
			nil,
		)
		if err != nil {
			return err
		}
		err = setRequestHeader(req)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		var response InstallationGetResponse
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
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
		accessToken, err := postOauthAccessToken(code)
		if err != nil {
			return err
		}
		githubResp, err := postGithubUser(accessToken)
		if err != nil {
			return err
		}
		query = `
			INSERT INTO user_githubs(user_id, github_id, github_name, created_at, updated_at)
				VALUES(:userId, :githubId, :githubName, :createdAt, :updatedAt)
			ON CONFLICT(user_id) 
			DO UPDATE SET
				github_id = :githubId,
				github_name = :githubName,
				updated_at = :updatedAt
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"userId":     userId,
			"githubId":   githubResp.Id,
			"githubName": githubResp.Login,
			"createdAt":  time.Now(),
			"updatedAt":  time.Now(),
		})
		return err
	})
}
