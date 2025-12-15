package query

import (
	"context"
	"maps"
	"slices"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type ReviewLoadStatusInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}

type ReviewLoadStatusListOutput struct {
	RepositoryId   int64
	RepositoryName string
	Total          int
	Status         []ReviewLoadStatusListOutputItem
}

type ReviewLoadStatusListOutputItem struct {
	GithubId   int64
	GithubName string
	GithubUrl  string
	Assigned   int
}

func (q *GithubQuery) GetReviewLoadStatus(ctx context.Context, input ReviewLoadStatusInput) ([]ReviewLoadStatusListOutput, error) {
	// オプションが設定されているかチェック
	options, err := q.optionRepository.FindsByCalendarId(ctx, input.CalendarId)
	if err != nil {
		return nil, err
	}
	hasOption := false
	for _, opt := range options {
		if opt.Id == option.OPTION_REVIEW_LOAD_STATUS {
			hasOption = true
			break
		}
	}
	if !hasOption {
		return nil, apperr.ForbiddenError("option not enabled")
	}
	// プルリクエストをGithubから取得
	prList, err := q.githubRepository.GetPullRequests(
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
	outputs := make([]ReviewLoadStatusListOutput, len(prList))
	for i, output := range prList {
		total := 0
		statusMap := map[int64]ReviewLoadStatusListOutputItem{}
		outputs[i].RepositoryId = output.Repository.Id
		outputs[i].RepositoryName = output.Repository.Name
		for _, pullRequest := range output.PullRequests {
			// レビュー依頼されていないものはスキップ
			if len(pullRequest.Reviewers) == 0 {
				continue
			}
			total += 1
			// 全てのレビュアーをチェック
			for _, reviewer := range pullRequest.Reviewers {
				status, ok := statusMap[reviewer.Id]
				if !ok {
					// マップにない場合新規作成
					status = ReviewLoadStatusListOutputItem{
						GithubId:   reviewer.Id,
						GithubName: reviewer.Name,
						GithubUrl:  reviewer.Url,
					}
				}
				// アサイン数+1
				status.Assigned += 1
				statusMap[reviewer.Id] = status
			}
		}
		outputs[i].Total = total
		outputs[i].Status = slices.Collect(maps.Values(statusMap))
	}
	return outputs, nil
}
