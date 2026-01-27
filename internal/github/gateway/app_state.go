package gateway

import (
	"context"
	"encoding/json"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AppStateModel struct {
	UserId     uuid.UUID `json:"userId"`
	CalendarId uuid.UUID `json:"calendarId"`
}

func (r *StateRedisRepository) GetAppState(
	ctx context.Context,
	state string,
) (uuid.UUID, uuid.UUID, error) {
	result, err := r.redisClient.GetDel(ctx, getAppStateKey(state)).Result()
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	var model AppStateModel
	if err := json.Unmarshal([]byte(result), &model); err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	return model.UserId, model.CalendarId, nil
}

func (r *StateRedisRepository) SaveAppState(
	ctx context.Context,
	userId, calendarId uuid.UUID,
	state string,
) error {
	return db.RunInTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		var isMember bool
		query := `
			SELECT 1 
				FROM calendar_members
			WHERE calendar_members.calendar_id = $2
				AND calendar_members.user_id = $1
				AND calendar_members.joined_at IS NOT NULL
		`
		err := tx.GetContext(ctx, &isMember, query, userId, calendarId)
		if err != nil {
			return err
		}
		if !isMember {
			return apperr.ForbiddenError("invalid member")
		}
		exp := time.Duration(time.Minute * 10)
		model := AppStateModel{
			UserId:     userId,
			CalendarId: calendarId,
		}
		json, err := json.Marshal(&model)
		if err != nil {
			return err
		}
		status := r.redisClient.SetEx(ctx, getAppStateKey(state), string(json), exp)
		return status.Err()
	})
}

func getAppStateKey(state string) string {
	return "github:apps:state:" + state
}
