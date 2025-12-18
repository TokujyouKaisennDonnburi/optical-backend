package psql

import (
	"context"
	"database/sql"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserModel struct {
	Id        uuid.UUID    `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Password  []byte       `db:"password_hash"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func FindUserById(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*user.User, error) {
	query := `
		SELECT 
			id, name, email, password_hash, created_at, updated_at, deleted_at
		FROM users
		WHERE 
			id = $1
	`
	userModel := UserModel{}
	err := tx.Get(&userModel, query, id)
	if err != nil {
		return nil, err
	}
	if userModel.DeletedAt.Valid {
		return nil, apperr.ForbiddenError("user is deleted")
	}
	return &user.User{
		Id:        userModel.Id,
		Name:      userModel.Name,
		Email:     user.Email(userModel.Email),
		Password:  userModel.Password,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
		DeletedAt: userModel.DeletedAt.Time,
	}, nil
}

func FindUserByEmail(ctx context.Context, tx *sqlx.Tx, email string) (*user.User, error) {
	query := `
		SELECT 
			id, name, email, password_hash, created_at, updated_at, deleted_at
		FROM users
		WHERE 
			users.email = $1
	`
	userModel := UserModel{}
	err := tx.GetContext(ctx, &userModel, query, email)
	if err != nil {
		return nil, err
	}
	if userModel.DeletedAt.Valid {
		return nil, apperr.ForbiddenError("user is deleted")
	}
	return &user.User{
		Id:        userModel.Id,
		Name:      userModel.Name,
		Email:     user.Email(userModel.Email),
		Password:  userModel.Password,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
		DeletedAt: userModel.DeletedAt.Time,
	}, nil
}

func FindUserByGithubSSO(ctx context.Context, tx *sqlx.Tx, githubId int64) (*user.User, error) {
	query := `
		SELECT 
			users.id, users.name, users.email, users.password_hash, users.created_at, users.updated_at, users.deleted_at
		FROM users
		JOIN user_githubs
			ON users.id = user_githubs.user_id
		WHERE 
			user_githubs.github_id = $1
			AND user_githubs.sso_login = true 
	`
	userModel := UserModel{}
	err := tx.GetContext(ctx, &userModel, query, githubId)
	if err != nil {
		return nil, err
	}
	return &user.User{
		Id:        userModel.Id,
		Name:      userModel.Name,
		Email:     user.Email(userModel.Email),
		Password:  userModel.Password,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
		DeletedAt: userModel.DeletedAt.Time,
	}, nil
}

func FindUsersByEmails(ctx context.Context, tx *sqlx.Tx, emails []string) ([]user.User, error) {
	query := `
		SELECT 
			id, name, email, password_hash, created_at, updated_at, deleted_at
		FROM users
		WHERE 
			users.email = ANY($1)
		ORDER BY users.id
	`
	userModels := []UserModel{}
	err := tx.SelectContext(ctx, &userModels, query, pq.Array(emails))
	if err != nil {
		return nil, err
	}
	users := make([]user.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = user.User{
			Id:        userModel.Id,
			Name:      userModel.Name,
			Email:     user.Email(userModel.Email),
			Password:  userModel.Password,
			CreatedAt: userModel.CreatedAt,
			UpdatedAt: userModel.UpdatedAt,
			DeletedAt: userModel.DeletedAt.Time,
		}
	}
	return users, nil
}
