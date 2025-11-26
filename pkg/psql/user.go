package psql

import (
	"context"
	"database/sql"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
	return &user.User{
		Id:        userModel.Id,
		Name:      userModel.Name,
		Email:     userModel.Email,
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
	err := tx.Get(&userModel, query, email)
	if err != nil {
		return nil, err
	}
	return &user.User{
		Id:        userModel.Id,
		Name:      userModel.Name,
		Email:     userModel.Email,
		Password:  userModel.Password,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
		DeletedAt: userModel.DeletedAt.Time,
	}, nil
}
