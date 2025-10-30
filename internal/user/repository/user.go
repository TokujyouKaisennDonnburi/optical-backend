package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
)

type UserRepository interface {
	Create (ctx context.Context, user *user.User) error
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

