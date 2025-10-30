package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
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
	query := `
		INSERT INTO users(id, name, email, password_hash, created_at, updated_at)
		VALUES(:id, :name, :email, :password, :createdAt, :updatedAt)
	`
	_, err := r.db.NamedExecContext(ctx, query, map[string]any{
		"id": user.Id,
		"name": user.Name,
		"email": user.Email,
		"password": user.Password,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
	if err != nil {
		return err
	}
	return nil
}
