package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (r *TodoPsqlRepository) AddItem(
	ctx context.Context,
	listId uuid.UUID,
	todoItem *todo.Item,
) error {
	return db.RunInTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		todoList, err := psql.FindTodoListById(ctx, tx, listId)
		if err != nil {
			return err
		}
		exists, err := psql.IsUserInCalendarMembers(ctx, tx, todoItem.UserId, todoList.CalendarId)
		if err != nil {
			return err
		}
		if !exists {
			return apperr.ForbiddenError("the user is not in calendar members")
		}
		query := `
			INSERT INTO todo_items(id, todo_list_id, user_id, name, is_done, created_at, updated_at)
			VALUES(:id, :todoListId, :userId, :name, :isDone, :createdAt, :updatedAt)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":         todoItem.Id,
			"todoListId": listId,
			"userId":     todoItem.UserId,
			"name":       todoItem.Name,
			"isDone":     todoItem.IsDone,
			"createdAt":  time.Now().UTC(),
			"updatedAt":  time.Now().UTC(),
		})
		return err
	})
}
