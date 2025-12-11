package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type GithubCalendarInstallInput struct {
	CalendarId     string
	InstallationId string
}

func (c *GithubCommand) InstallToCalendar(ctx context.Context, input GithubCalendarInstallInput) error {
	if input.InstallationId == "" {
		return apperr.ValidationError("invalid installationId")
	}
	calendarId, err := uuid.Parse(input.CalendarId)
	if err != nil {
		return err
	}
	return c.githubRepository.InstallToCalendar(ctx, calendarId, input.InstallationId)
}
