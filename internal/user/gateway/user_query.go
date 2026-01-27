package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (r *UserPsqlRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var err error
	var user *user.User
	err = db.RunInTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		user, err = psql.FindUserByEmail(ctx, tx, email)
		return err
	})
	return user, err
}

func (r *UserPsqlRepository) FindById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var err error
	var user *user.User
	err = db.RunInTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		user, err = psql.FindUserById(ctx, tx, id)
		return err
	})
	return user, err
}
