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
	githubUser, err := api.GetGithubUser(accessToken)
	if err != nil {
		return nil, err
	}
	password, err := security.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	var newUser *user.User
	err = db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		newUser, err = psql.FindUserByGithubSSO(ctx, tx, githubUser.Id)
		if err == nil {
			return updateGithubUser(ctx, tx, newUser.Id, githubUser)
		} else {
			newUser, err = user.NewUser(githubUser.Name, githubUser.Email, password)
			if err != nil {
				return err
			}
			avatar, err := user.NewAvatar(githubUser.AvatarUrl, false)
			if err != nil {
				return err
			}
			return createGithubUser(ctx, tx, newUser, avatar, githubUser)
		}
	})
	return newUser, err
}

func createGithubUser(ctx context.Context, tx *sqlx.Tx, newUser *user.User, avatar *user.Avatar, githubUser *github.User) error {
	query := `
		INSERT INTO users(id, name, email, password_hash, created_at, updated_at)
		VALUES(:id, :name, :email, :password, :createdAt, :updatedAt)
	`
	_, err := tx.NamedExecContext(ctx, query, map[string]any{
		"id":        newUser.Id,
		"name":      newUser.Name,
		"email":     newUser.Email,
		"password":  newUser.Password,
		"createdAt": time.Now().UTC(),
		"updatedAt": time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	query = `
		INSERT INTO user_githubs(user_id, github_id, github_name, github_email, sso_login, created_at, updated_at)
		VALUES(:userId, :githubId, :githubName, :githubEmail, true, :createdAt, :updatedAt)
	`
	_, err = tx.NamedExecContext(ctx, query, map[string]any{
		"userId":      newUser.Id,
		"githubId":    githubUser.Id,
		"githubName":  githubUser.Name,
		"githubEmail": githubUser.Email,
		"createdAt":   time.Now().UTC(),
		"updatedAt":   time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	query = `
		INSERT INTO avatars(id, url, is_relative_path)
		VALUES(:id, :url, :isRelativePath)
	`
	_, err = tx.NamedExecContext(ctx, query, map[string]any{
		"id":             avatar.Id,
		"url":            avatar.Url,
		"isRelativePath": avatar.IsRelativePath,
	})
	if err != nil {
		return err
	}
	query = `
		INSERT INTO user_profiles(user_id, avatar_id)	
		VALUES(:userId, :avatarId)
	`
	_, err = tx.NamedExecContext(ctx, query, map[string]any{
		"userId":   newUser.Id,
		"avatarId": avatar.Id,
	})
	if err != nil {
		return err
	}
	return err
}

func updateGithubUser(ctx context.Context, tx *sqlx.Tx, userId uuid.UUID, githubUser *github.User) error {
	query := `
		INSERT INTO user_githubs(user_id, github_id, github_name, github_email, sso_login, created_at, updated_at)
			VALUES(:userId, :githubId, :githubName, :githubEmail, true, :createdAt, :updatedAt)
		ON CONFLICT(user_id) 
		DO UPDATE SET
			github_id = :githubId,
			github_name = :githubName,
			github_email = :githubEmail,
			updated_at = :updatedAt
	`
	_, err := tx.NamedExecContext(ctx, query, map[string]any{
		"userId":      userId,
		"githubId":    githubUser.Id,
		"githubName":  githubUser.Name,
		"githubEmail": githubUser.Email,
		"createdAt":   time.Now().UTC(),
		"updatedAt":   time.Now().UTC(),
	})
	return err
}
