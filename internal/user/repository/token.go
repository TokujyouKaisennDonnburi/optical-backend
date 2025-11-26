package repository

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type TokenRepository interface {
	AddToWhitelist(refreshToken *user.RefreshToken) error
	IsWhitelisted(tokenId uuid.UUID) error
}
