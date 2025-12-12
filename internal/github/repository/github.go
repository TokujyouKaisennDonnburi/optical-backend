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
	LinkUser(
		ctx context.Context,
		userId uuid.UUID,
		code string,
	) error
}
