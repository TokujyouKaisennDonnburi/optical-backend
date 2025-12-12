package repository

import (
	"context"

	"github.com/google/uuid"
)

type GithubRepository interface {
	InstallToCalendar(
		ctx context.Context,
		userId, calendarId uuid.UUID,
		code, installationId string,
	) error
	LinkUser(
		ctx context.Context,
		userId uuid.UUID,
		code string,
	) error
}
