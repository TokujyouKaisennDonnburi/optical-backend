package gateway

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query/output"
	"github.com/google/uuid"
)

// gateway用構造体
type IsLinkedUserModel struct {
	GithubId    string    `db:"github_id"`
	GithubName  string    `db:"github_name"`
	GithubEmail string    `db:"github_email"`
	SsoLogin    bool      `db:"sso_login"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// GitHubアカウントが連携されているか取得
func (r *GithubApiRepository) IsLinkedUser(
	ctx context.Context,
	userId uuid.UUID,
) (
	*output.IsLinkedUserQueryOutput,
	error,
) {
	query := `
		SELECT
		      github_id, github_name, github_email, sso_login, updated_at
		FROM user_githubs
		WHERE user_id = $1
		`

	var model IsLinkedUserModel
	err := r.db.GetContext(ctx, &model, query, userId)

	// DBに無ければ未連携
	if errors.Is(err, sql.ErrNoRows) {
		return &output.IsLinkedUserQueryOutput{IsLinked: false}, nil
	}

	if err != nil {
		return nil, err
	}

	return &output.IsLinkedUserQueryOutput{
		IsLinked:    true,
		GithubId:    model.GithubId,
		GithubName:  model.GithubName,
		GithubEmail: model.GithubEmail,
		IsSsoLogin:  model.SsoLogin,
		LinkedAt:    model.UpdatedAt,
	}, nil
}
