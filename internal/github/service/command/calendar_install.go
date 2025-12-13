package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
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
	// キャッシュ処理を行い失敗した場合エラーのみ出力
	cacheFunc := func() error {
		organization, err := api.GetInstalledOrganization(ctx, input.InstallationId)
		if err != nil {
			return err
		}
		now := time.Now()
		if organization.TokenExpiresAt.Before(now) {
			return errors.New("cache error: invalid time")
		}
		return c.stateRepository.SaveOrganization(ctx, organization, organization.TokenExpiresAt.Sub(now))
	}
	if err := cacheFunc(); err != nil {
		fmt.Printf("servive@cache error: %s\n", err.Error())
	}
	return c.githubRepository.InstallToCalendar(ctx, userId, calendarId, input.Code, input.InstallationId)
}
