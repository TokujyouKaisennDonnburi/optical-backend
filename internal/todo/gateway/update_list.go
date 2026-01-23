package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (r *TodoPsqlRepository) UpdateList(
	ctx context.Context,
	userId, listId uuid.UUID,
	updateFn func(*todo.List) (*todo.List, error),
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		todoList, err := psql.FindTodoListByIdAndUserId(ctx, tx, listId, userId)
		if err != nil {
			return err
		}
		// 更新処理
		todoList, err = updateFn(todoList)
		if err != nil {
			return err
		}
		// DB更新
		query := `
			UPDATE todo_lists SET
				name = :name,
				updated_at = :updatedAt
			WHERE
				id = :id
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":         todoList.Id,
			"name":       todoList.Name,
			"updatedAt":  time.Now().UTC(),
		})
		return err
	})
}
