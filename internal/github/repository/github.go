package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type GithubRepository interface {
	InstallToCalendar(
		ctx context.Context,
		userId, calendarId uuid.UUID,
		code, installationId string,
	) error
	CreateUser(
		ctx context.Context,
		code string,
	) (*user.User, error)
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
	GetMilestones(
		ctx context.Context,
		userId, calendarId uuid.UUID,
		getFn func(installationId string) (*github.Organization, error),
	) (
		[]github.Milestones,
		error,
	)

	// GitHubアカウントが連携されているか
	IsLinkedUser(
		ctx context.Context,
		userId uuid.UUID,
	) (
		*output.IsLinkedUserQueryOutput,
		error,
	)
}
