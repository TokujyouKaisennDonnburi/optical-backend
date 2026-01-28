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

// 複数のメールアドレスから登録ユーザーがいるか確認
func (r *UserPsqlRepository) FindByEmails(ctx context.Context, emails []string) ([]*user.User, error) {
	var users []user.User
	err := db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		var err error
		users, err = psql.FindUsersByEmails(ctx, tx, emails)
		return err
	})
	if err != nil {
		return nil, err
	}
	result := make([]*user.User, len(users))
	for i := range users {
		result[i] = &users[i]
	}
	return result, nil
}
