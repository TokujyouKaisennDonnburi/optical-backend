package gateway

import (
	"context"
	"errors"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserPsqlRepository struct {
	db *sqlx.DB
}

func NewUserPsqlRepository(db *sqlx.DB) *UserPsqlRepository {
	if db == nil {
		panic("db is nil")
	}
	return &UserPsqlRepository{
		db: db,
	}
}

func (r *UserPsqlRepository) Create(ctx context.Context, user *user.User) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		query := `
			INSERT INTO users(id, name, email, password_hash, created_at, updated_at)
			VALUES(:id, :name, :email, :password, :createdAt, :updatedAt)
		`
		_, err := tx.NamedExecContext(ctx, query, map[string]any{
			"id":        user.Id,
			"name":      user.Name,
			"email":     user.Email,
			"password":  user.Password,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		})
		if err != nil {
			return err
		}
		return err
	})
}

func (r *UserPsqlRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	updateFn func(*user.User) error,
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		user, err := psql.FindUserById(ctx, tx, id)
		if err != nil {
			return err
		}
		err = updateFn(user)
		if err != nil {
			return err
		}
		query := `
			UPDATE users SET
				name = :name,
				email = :email,
				updated_at = :updatedAt
			WHERE
				users.id = :id
		`
		result, err := tx.NamedExecContext(ctx, query, map[string]any{
			"id":        user.Id,
			"name":      user.Name,
			"email":     user.Email.String(),
			"updatedAt": time.Now().UTC(),
		})
		if err != nil {
			return err
		}
		if rows, _ := result.RowsAffected(); rows == 0 {
			return errors.New("failed to update user")
		}
		return nil
	})
}
