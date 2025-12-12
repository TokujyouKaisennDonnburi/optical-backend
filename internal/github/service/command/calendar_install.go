package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
)

type GithubCalendarInstallInput struct {
	State          string
	InstallationId string
}

func (c *GithubCommand) InstallToCalendar(ctx context.Context, input GithubCalendarInstallInput) error {
	if input.State == "" {
		return apperr.ValidationError("invalid state")
	}
	if input.InstallationId == "" {
		return apperr.ValidationError("invalid installationId")
	}
	calendarId, err := c.stateRepository.GetAppState(ctx, input.State)
	if err != nil {
		return err
	}
	return c.githubRepository.InstallToCalendar(ctx, calendarId, input.InstallationId)
}
