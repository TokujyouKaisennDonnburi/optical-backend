package gateway

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
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

type UserModel struct {
	Id        uuid.UUID    `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Password  []byte       `db:"password_hash"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
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
	query := `
		SELECT 
			id, name, email, password_hash, created_at, updated_at, deleted_at
		FROM users
		WHERE 
			users.email = $1
	`
	userModel := UserModel{}
	err := r.db.Get(&userModel, query, email)
	if err != nil {
		return nil, errors.New(err.Error() + ", email:" + email)
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

func (r *UserPsqlRepository) FindById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
		SELECT 
			id, name, email, password_hash, created_at, updated_at, deleted_at
		FROM users
		WHERE 
			id = $1
	`
	userModel := UserModel{}
	err := r.db.Get(&userModel, query, id)
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
