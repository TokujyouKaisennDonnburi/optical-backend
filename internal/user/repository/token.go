package repository

import "github.com/TokujouKaisenDonburi/optical-backend/internal/user"

type TokenRepository interface {
	AddToWhitelist(refreshToken *user.RefreshToken) error
}
