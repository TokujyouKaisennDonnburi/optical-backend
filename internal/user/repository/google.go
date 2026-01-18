package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
)

type GoogleRepository interface {
	GetTokenByCode(ctx context.Context, code string) (string, error)
	GetUserByToken(ctx context.Context, token string) (*user.GoogleUser, error)
	CreateUser(
		ctx context.Context,
		user *user.User,
		avatar *user.Avatar,
		googleUser *user.GoogleUser,
	) error
}
