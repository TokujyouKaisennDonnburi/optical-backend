package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type TokenRepository interface {
	AddToWhitelist(ctx context.Context, refreshToken *user.RefreshToken) error
	IsWhitelisted(ctx context.Context, userId uuid.UUID, tokenId uuid.UUID) error
}
