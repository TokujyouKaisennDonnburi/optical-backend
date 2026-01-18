package psql

import (
	"context"
	"database/sql"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/security"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InstallationIdAndGithubIdModel struct {
	GithubId       int64  `db:"github_id"`
	InstallationId string `db:"installation_id"`
}

// カレンダーIDとユーザーIDからインストールIDを取得
func FindInstallationIdAndGithubId(
	ctx context.Context,
	tx *sqlx.Tx,
	userId, calendarId uuid.UUID,
	installationIdEncryptionKey []byte,
) (int64, string, error) {
	var model InstallationIdAndGithubIdModel
	query := `
		SELECT installation_id, user_githubs.github_id
		FROM (
			SELECT * FROM calendar_githubs
			WHERE calendar_id = $2
		) c_githubs
		JOIN calendar_members
			ON calendar_members.calendar_id = c_githubs.calendar_id
		JOIN user_githubs
			ON user_githubs.user_id = $1
	`
	err := tx.GetContext(ctx, &model, query, userId, calendarId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", apperr.ForbiddenError("no github installation or no calendar access")
		}
		return 0, "", err
	}
	// 暗号化されたinstallation_idを復号化
	decryptedInstallationId, err := security.Decrypt(model.InstallationId, installationIdEncryptionKey)
	if err != nil {
		return 0, "", apperr.InternalServerError("failed to decrypt installation_id")
	}
	return model.GithubId, decryptedInstallationId, nil
}
