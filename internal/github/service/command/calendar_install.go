package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
)

type GithubCalendarInstallInput struct {
	Code           string
	State          string
	InstallationId string
}

func (c *GithubCommand) InstallToCalendar(ctx context.Context, input GithubCalendarInstallInput) error {
	if input.Code == "" {
		return apperr.ValidationError("invalid code")
	}
	if input.State == "" {
		return apperr.ValidationError("invalid state")
	}
	if input.InstallationId == "" {
		return apperr.ValidationError("invalid installationId")
	}
	userId, calendarId, err := c.stateRepository.GetAppState(ctx, input.State)
	if err != nil {
		return err
	}
	return c.githubRepository.InstallToCalendar(ctx, userId, calendarId, input.Code, input.InstallationId)
}
