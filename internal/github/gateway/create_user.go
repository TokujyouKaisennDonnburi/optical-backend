package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/security"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (r *GithubApiRepository) CreateUser(
	ctx context.Context,
	code string,
) (*user.User, error) {
	accessToken, err := api.PostOauthAccessToken(code)
	if err != nil {
		return nil, err
	}
	githubResp, err := api.GetGithubUser(accessToken)
	if err != nil {
		return nil, err
	}
	email, err := api.GetGithubPrimaryEmail(accessToken)
	if err != nil {
		return nil, err
	}
	password, err := security.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	githubUser := &github.User{
		Id:        githubResp.Id,
		Login:     githubResp.Login,
		Email:     email,
		Url:       githubResp.Url,
		AvatarUrl: githubResp.AvatarUrl,
	}
	var newUser *user.User
	err = db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		newUser, err = psql.FindUserByGithubId(ctx, tx, githubUser.Id)
		if err == nil {
			return updateGithubUser(ctx, tx, newUser.Id, githubUser)
		} else {
			newUser, err = user.NewUser(githubUser.Login, githubUser.Email, password)
			if err != nil {
				return err
			}
			return createGithubUser(ctx, tx, newUser, githubUser)
		}
	})
	return newUser, err
}

func createGithubUser(ctx context.Context, tx *sqlx.Tx, newUser *user.User, githubUser *github.User) error {
	query := `
			INSERT INTO users(id, name, email, password_hash, created_at, updated_at)
			VALUES(:id, :name, :email, :password, :createdAt, :updatedAt)
		`
	_, err := tx.NamedExecContext(ctx, query, map[string]any{
		"id":        newUser.Id,
		"name":      newUser.Name,
		"email":     newUser.Email,
		"password":  newUser.Password,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	})
	if err != nil {
		return err
	}
	query = `
			INSERT INTO user_profiles(user_id, image_url)	
			VALUES(:userId, :imageUrl)
		`
	_, err = tx.NamedExecContext(ctx, query, map[string]any{
		"userId":   newUser.Id,
		"imageUrl": githubUser.AvatarUrl,
	})
	if err != nil {
		return err
	}
	query = `
			INSERT INTO user_githubs(user_id, github_id, github_name, github_email, created_at, updated_at)
			VALUES(:userId, :githubId, :githubName, :githubEmail, :createdAt, :updatedAt)
		`
	_, err = tx.NamedExecContext(ctx, query, map[string]any{
		"userId":     newUser.Id,
		"githubId":   githubUser.Id,
		"githubName": githubUser.Login,
		"githubEmail": githubUser.Email,
		"createdAt":  time.Now(),
		"updatedAt":  time.Now(),
	})
	return err
}

func updateGithubUser(ctx context.Context, tx *sqlx.Tx, userId uuid.UUID, githubUser *github.User) error {
	query := `
			INSERT INTO user_githubs(user_id, github_id, github_name, github_email, created_at, updated_at)
			VALUES(:userId, :githubId, :githubName, :githubEmail, :createdAt, :updatedAt)
			ON CONFLICT(user_id) 
			DO UPDATE SET
				github_id = :githubId,
				github_name = :githubName,
				github_email = :githubEmail,
				updated_at = :updatedAt
		`
	_, err := tx.NamedExecContext(ctx, query, map[string]any{
		"userId":     userId,
		"githubId":   githubUser.Id,
		"githubName": githubUser.Login,
		"githubEmail": githubUser.Email,
		"createdAt":  time.Now(),
		"updatedAt":  time.Now(),
	})
	return err
}
