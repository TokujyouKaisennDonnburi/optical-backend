package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (r *GithubApiRepository) GetPullRequests(
	ctx context.Context,
	userId, calendarId uuid.UUID,
	getFn func(installationId string) (*github.Organization, error),
) (
	[]output.GithubPullRequestListQueryOutput,
	error,
) {
	var outputs []output.GithubPullRequestListQueryOutput
	err := db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		githubId, installationId, err := psql.FindInstallationIdAndGithubId(ctx, tx, userId, calendarId)
		if err != nil {
			return err
		}
		organization, err := getFn(installationId)
		if err != nil {
			return err
		}
		fmt.Println("Organization Repositories", organization.Repositories)
		for _, repos := range organization.Repositories {
			prList, err := getPullRequests(ctx, organization.AccessToken, organization.Login, repos.Name)
			if err != nil {
				fmt.Printf("repository error: %s\n", err.Error())
				continue
			}
			outputs = append(outputs, output.GithubPullRequestListQueryOutput{
				GithubId:     githubId,
				Repository:   repos,
				PullRequests: prList,
			})
		}
		return nil
	})
	return outputs, err
}

func getPullRequests(ctx context.Context, accessToken, owner, repos string) ([]github.PullRequest, error) {
	client := http.Client{}
	requestUrl := api.GITHUB_BASE_URL+"/repos/"+owner+"/"+repos+"/pulls"
	requestUrl += "?per_page=100"
	// インストール済みアプリ取得リクエスト
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		requestUrl,
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("X-GitHub-Api-Version", api.GITHUB_API_VERSION)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get pull requests: %d", resp.StatusCode)
	}
	// WARNING: ドメインから直接参照
	var pullRequests []github.PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
		return nil, err
	}
	fmt.Println("PullRequestsResponse:", pullRequests)
	return pullRequests, nil
}
