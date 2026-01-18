package repository

import (
	"context"

	"github.com/google/uuid"
)

type OauthStateRepository interface {
	GetClientId() string
	GetRedirectUri() string
	SaveOauthState(
		ctx context.Context,
		userId uuid.UUID,
		state string,
	) error
	GetOauthState(
		ctx context.Context,
		state string,
	) (uuid.UUID, error)
}
