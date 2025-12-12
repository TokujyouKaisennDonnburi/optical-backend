package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
	"github.com/google/uuid"
)

type GithubReviewRequestsQueryInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}

func (q *GithubQuery) GetReviewRequests(ctx context.Context, input GithubReviewRequestsQueryInput) ([]github.PullRequest, error) {
	outputs, err := q.githubRepository.GetPullRequests(
		ctx,
		input.UserId,
		input.CalendarId,
		func(installationId string) (*github.Organization, error) {
			organization, err := q.stateRepository.GetOrganization(ctx, installationId)
			if err != nil {
				return nil, err
			}
			repositories, err := api.GetInstalledRepositories(ctx, organization.AccessToken)
			if err != nil {
				return nil, err
			}
			organization.SetRepositories(repositories)
			return organization, nil
		},
	)
	if err != nil {
		return nil, err
	}
	var pullRequests []github.PullRequest
	for _, output := range outputs {
		for _, pullRequest := range output.PullRequests {
			// レビューされていないものはスキップ
			if pullRequest.ReviewComments != 0 {
				continue
			}
			// クローズしているものはスキップ
			if pullRequest.State == github.PULL_REQUEST_STATE_CLOSE {
				continue
			}
			for _, reviewer := range pullRequest.Reviewers {
				if reviewer.Id != output.GithubId {
					continue
				}
				// アサインされている場合リストに追加
				pullRequests = append(pullRequests, pullRequest)
				break
			}
		}
	}
	return pullRequests, nil
}
