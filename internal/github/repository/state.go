package repository

import (
	"context"

	"github.com/google/uuid"
)

type StateRepository interface {
	SaveAppState(
		ctx context.Context,
		userId, calendarId uuid.UUID,
		state string,
	) error
	GetAppState(
		ctx context.Context,
		state string,
	) (uuid.UUID, uuid.UUID, error)
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
