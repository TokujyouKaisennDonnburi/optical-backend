package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create (ctx context.Context, user *user.User) error
	FindById(ctx context.Context, id uuid.UUID) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

