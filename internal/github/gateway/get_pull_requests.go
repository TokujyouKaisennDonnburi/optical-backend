package gateway

import (
	"context"
	"fmt"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// プルリクエストを取得
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
		// データベースからインストールIDを取得
		githubId, installationId, err := psql.FindInstallationIdAndGithubId(ctx, tx, userId, calendarId)
		if err != nil {
			return err
		}
		// インストールIDを元に組織を取得
		organization, err := getFn(installationId)
		if err != nil {
			return err
		}
		for _, repos := range organization.Repositories {
			// APIからそれぞれのリポジトリのプルリクエストを取得
			prList, err := api.GetPullRequests(ctx, organization.AccessToken, organization.Login, repos.Name)
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
