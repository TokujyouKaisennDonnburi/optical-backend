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

func (r *TodoPsqlRepository) UpdateItem(
	ctx context.Context,
	userId, itemId uuid.UUID,
	updateFn func(*todo.Item) (*todo.Item, error),
) error {
	return db.RunInTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		todoItem, err := psql.FindTodoItemById(ctx, tx, itemId)
		if err != nil {
			return err
		}
		exists, err := psql.IsUserInTodoListCalendarMembers(ctx, tx, userId, todoItem.ListId)
		if err != nil {
			return err
		}
		if !exists {
			return apperr.ForbiddenError("the user is not in calendar members")
		}
		// 更新処理
		todoItem, err = updateFn(todoItem)
		if err != nil {
			return err
		}
		// DB更新
		query := `
			UPDATE todo_items SET
				name = :name,
				is_done = :isDone,
				updated_at = :updatedAt
			WHERE
				id = :id
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":         todoItem.Id,
			"name":       todoItem.Name,
			"isDone":     todoItem.IsDone,
			"updatedAt":  time.Now().UTC(),
		})
		return err
	})
}
