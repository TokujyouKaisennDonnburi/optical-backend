package gateway

import (
	"context"
	"database/sql"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/query/output"
	"github.com/google/uuid"
)

type UserProfileModel struct {
	Id        uuid.UUID      `db:"id"`
	Name      string         `db:"name"`
	Email     string         `db:"email"`
	ImageUrl  sql.NullString `db:"image_url"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func (r *UserPsqlRepository) FindProfileById(ctx context.Context, id uuid.UUID) (*output.UserQueryOutput, error) {
	query := `
		SELECT 
			id, name, email, created_at, updated_at, user_profiles.image_url
		FROM users
		LEFT JOIN user_profiles
			ON users.id = user_profiles.user_id
		WHERE 
			users.id = $1
			AND users.deleted_at IS NULL
	`
	var model UserProfileModel
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		return nil, err
	}
	return &output.UserQueryOutput{
		Id:    model.Id,
		Name:  model.Name,
		Email: model.Email,
		Avatar: user.Avatar{
			Url:   model.ImageUrl.String,
			Valid: model.ImageUrl.Valid,
		},
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}
