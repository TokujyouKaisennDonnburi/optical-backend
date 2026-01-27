package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func (r *UserPsqlRepository) FindsByEmails(
	ctx context.Context,
	emails []string,
) ([]user.User, error) {
	var userModels []psql.UserModel
	err := db.RunInTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		query := `
			SELECT 
				id, name, email, password_hash, created_at, updated_at, deleted_at
			FROM users
			WHERE 
				users.email = ANY($1)
			ORDER BY users.id
		`
		err := tx.SelectContext(ctx, &userModels, query, pq.Array(emails))
		return err
	})
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
