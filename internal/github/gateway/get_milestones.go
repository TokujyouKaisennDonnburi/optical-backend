package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// マイルストーンを取得
func (r *GithubApiRepository) GetMilestones(
	ctx context.Context,
	userId, calendarId uuid.UUID,
	getFn func(installationId string) (*github.Organization, error),
) (
	[]github.Milestones,
	error,
) {
	var milestones []github.Milestones
	err := db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		_, installationId, err := psql.FindInstallationIdAndGithubId(ctx, tx, userId, calendarId)
		if err != nil {
			return err
		}
		organization, err := getFn(installationId)
		if err != nil {
			return err
		}
		for _, repos := range organization.Repositories {
			repoMilestones, err := api.GetMilestones(
				ctx,
				organization.AccessToken,
				organization.Login,
				repos.Name,
				github.MILESTONES_STATE_OPEN,
			)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"organization_id":    organization.Id,
					"organization_login": organization.Login,
					"repository_name":    repos.Name,
				}).WithError(err).Error("failed to get milestones")
				continue
			}
			milestones = append(milestones, repoMilestones...)
		}
		return nil
	})
	return milestones, err
}
