package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (r *StateRedisRepository) GetAppState(
	ctx context.Context,
	state string,
) (uuid.UUID, error) {
	result, err := r.redisClient.GetDel(ctx, getAppStateKey(state)).Result()
	if err != nil {
		return uuid.Nil, err
	}
	userId, err := uuid.Parse(result)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}

func (r *StateRedisRepository) SaveAppState(
	ctx context.Context,
	userId, calendarId uuid.UUID,
	state string,
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		var isMember bool
		query := `
			SELECT 1 FROM calendar_members
			WHERE calendar_members.calendar_id = $2
			AND calendar_members.user_id = $1
			AND calendar_members.joined_at IS NOT NULL
		`
		err := tx.GetContext(ctx, &isMember, query, userId, calendarId)
		if err != nil {
			return err
		}
		if !isMember {
			return apperr.ForbiddenError("")
		}
		exp := time.Duration(time.Minute * 10)
		status := r.redisClient.SetEx(ctx, getAppStateKey(state), calendarId.String(), exp)
		return status.Err()
	})
}

func getAppStateKey(state string) string {
	return "github:apps:state:" + state
}
