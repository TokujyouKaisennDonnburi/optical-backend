package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query/output"
	"github.com/google/uuid"
)

func (q *GithubQuery) IsLinkedUser(
	ctx context.Context,
	userId uuid.UUID,
) (*output.IsLinkedUserQueryOutput, error) {
	return q.githubRepository.IsLinkedUser(ctx, userId)
}
