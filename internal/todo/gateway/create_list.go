package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/jmoiron/sqlx"
)

func (r *TodoPsqlRepository) CreateList(
	ctx context.Context,
	list *todo.List,
) error {
	return db.RunInTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		exists, err := psql.IsUserInCalendarMembers(ctx, tx, list.UserId, list.CalendarId)
		if err != nil {
			return err
		}
		if !exists {
			return apperr.ForbiddenError("the user is not in calendar members")
		}
		query := `
			INSERT INTO todo_lists(id, user_id, calendar_id, name, created_at, updated_at)
			VALUES(:id, :userId, :calendarId, :name, :createdAt, :updatedAt)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":         list.Id,
			"userId":     list.UserId,
			"calendarId": list.CalendarId,
			"name":       list.Name,
			"createdAt":  time.Now().UTC(),
			"updatedAt":  time.Now().UTC(),
		})
		return err
	})
}

