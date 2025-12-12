package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query/output"
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
	GetPullRequests(
		ctx context.Context,
		userId, calendarId uuid.UUID,
		getFn func(installationId string) (*github.Organization, error),
	) (
		[]output.GithubPullRequestListQueryOutput,
		error,
	)
}
