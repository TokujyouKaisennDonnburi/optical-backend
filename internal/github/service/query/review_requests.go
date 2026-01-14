package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type GithubReviewRequestsQueryInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}

func (q *GithubQuery) GetReviewRequests(ctx context.Context, input GithubReviewRequestsQueryInput) ([]github.PullRequest, error) {
	// オプションが設定されているかチェック
	options, err := q.optionRepository.FindsByCalendarId(ctx, input.CalendarId)
	if err != nil {
		return nil, err
	}
	hasOption := false
	for _, opt := range options {
		if opt.Id == option.OPTION_PULL_REQUEST_REVIEW_WAIT_COUNT {
			hasOption = true
			break
		}
	}
	if !hasOption {
		return nil, apperr.ForbiddenError("option not enabled")
	}
	// プルリクエストをGithubから取得
	outputs, err := q.githubRepository.GetPullRequests(
		ctx,
		input.UserId,
		input.CalendarId,
		func(installationId string) (*github.Organization, error) {
			// 組織を取得
			organization, err := q.stateRepository.GetOrganization(ctx, installationId)
			if err != nil {
				return nil, err
			}
			// APIから最新のリポジトリを取得
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
			// ドラフトはスキップ
			if pullRequest.Draft {
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
