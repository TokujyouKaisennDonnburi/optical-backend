package repository

import (
	"context"

	"github.com/google/uuid"
)

type GithubRepository interface {
	InstallToCalendar(
		ctx context.Context,
		calendarId uuid.UUID,
		installationId string,
	) error
	SaveOauthState(
		ctx context.Context,
		userId uuid.UUID,
		state string,
	) error
	LinkUser(
		ctx context.Context,
		code, state string,
	) error
}
