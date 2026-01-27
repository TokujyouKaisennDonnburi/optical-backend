package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/service/query/output"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		user *user.User,
	) error
	Update(
		ctx context.Context,
		userId uuid.UUID,
		updateFn func (*user.User) error,
	) error
	FindById(ctx context.Context, id uuid.UUID) (*user.User, error)
	FindProfileById(ctx context.Context, id uuid.UUID) (*output.UserQueryOutput, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindsByEmails(ctx context.Context, email []string) ([]user.User, error)
}
