package gateway

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

// gateway用構造体
type IsInstalledGithubAppModel struct {
	GithubId       string `db:"github_id"`
	GithubName     string `db:"github_name"`
	InstallationId string `db:"installation_id"`
	UpdateAt       string `db:"updated_at"` // InstalledAtとする
}

// GithubAppが連携されているか確認
func (r *GithubApiRepository) IsInstalledGithubApp(
	ctx context.Context,
	userId, calendarId uuid.UUID,
) (
	*output.IsInstalledGithubAppQueryOutput,
	error,
) {
	query := `
		SELECT 
			COALESCE(cg.github_id, '') as github_id,
			COALESCE(cg.github_name, '') as github_name,
			COALESCE(cg.installation_id, '') as installation_id,
			COALESCE(cg.updated_at::text, '') as updated_at
		FROM calendar_members cm
		LEFT JOIN calendar_githubs cg ON cg.calendar_id = cm.calendar_id
		WHERE cm.calendar_id = $1 AND cm.user_id = $2
		`
	var model IsInstalledGithubAppModel
	err := r.db.GetContext(ctx, &model, query, calendarId, userId)

	// カレンダーメンバーでない
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.ForbiddenError("not a member of this calendar")
	}

	// その他エラー
	if err != nil {
		return nil, err
	}

	// メンバーだが未インストール
	if model.GithubId == "" {
		return &output.IsInstalledGithubAppQueryOutput{IsInstalled: false}, nil
	}

	return &output.IsInstalledGithubAppQueryOutput{
		IsInstalled:    true,
		GithubId:       model.GithubId,
		GithubName:     model.GithubName,
		InstallationId: model.InstallationId,
		InstalledAt:    model.UpdateAt,
	}, nil
}
