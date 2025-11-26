package gateway

import (
	"context"

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
	query := `
		INSERT INTO users(id, name, email, password_hash, created_at, updated_at)
		VALUES(:id, :name, :email, :password, :createdAt, :updatedAt)
	`
	_, err := r.db.NamedExecContext(ctx, query, map[string]any{
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
	return nil
}

func (r *UserPsqlRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var err error
	var user *user.User
	err = db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		user, err = psql.FindUserByEmail(ctx, tx, email)
		return err
	})
	return user, err
}

func (r *UserPsqlRepository) FindById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var err error
	var user *user.User
	err = db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		user, err = psql.FindUserById(ctx, tx, id)
		return err
	})
	return user, err
}
