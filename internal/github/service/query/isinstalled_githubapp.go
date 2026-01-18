package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query/output"
	"github.com/google/uuid"
)

func (q *GithubQuery) IsInstalledGithubApp(
	ctx context.Context,
	userid, calendarId uuid.UUID,
) (*output.IsInstalledGithubAppQueryOutput, error) {
	return q.githubRepository.IsInstalledGithubApp(ctx, userid, calendarId)
}
