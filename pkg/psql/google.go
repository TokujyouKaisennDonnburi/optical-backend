package psql

import (
	"context"
	"database/sql"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/jmoiron/sqlx"
)

func FindUserByGoogleId(
	ctx context.Context,
	tx *sqlx.Tx,
	googleId string,
) (*user.User, error) {
	query := `
		SELECT 
			users.id, users.name, users.email, users.password_hash, users.created_at, users.updated_at, users.deleted_at
		FROM users
		JOIN google_ids
			ON users.id = google_ids.user_id
		WHERE 
			google_ids.google_id = $1
	`
	userModel := UserModel{}
	err := tx.Get(&userModel, query, googleId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperr.NotFoundError(err.Error())
		}
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
