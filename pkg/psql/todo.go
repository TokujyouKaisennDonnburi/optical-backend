package psql

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TodoListAndItemModel struct {
	Id         uuid.UUID `db:"id"`
	UserId     uuid.UUID `db:"user_id"`
	CalendarId uuid.UUID `db:"calendar_id"`
	Name       string    `db:"name"`
	ItemId     uuid.UUID `db:"item_id"`
	ItemUserId uuid.UUID `db:"item_user_id"`
	ItemName   string    `db:"item_name"`
	ItemIsDone bool      `db:"item_is_done"`
}

func FindTodoListById(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*todo.List, error) {
	query := `
		SELECT id, user_id, calendar_id, name, created_at, updated_at
		FROM todo_lists
		JOIN todo_items
			ON todo_lists.id = todo_items.todo_list_id
		WHERE id = $1
	`
	var models []TodoListAndItemModel
	err := tx.SelectContext(ctx, &models, query, id)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, apperr.NotFoundError("todo list not found")
	}
	todoItems := make([]todo.Item, len(models))
	for i, model := range models {
		todoItems[i] = todo.Item{
			Id:     model.ItemId,
			UserId: model.ItemUserId,
			Name:   model.ItemName,
			IsDone: model.ItemIsDone,
		}
	}
	return &todo.List{
		Id:         models[0].Id,
		UserId:     models[0].UserId,
		CalendarId: models[0].CalendarId,
		Name:       models[0].Name,
		Items:      todoItems,
	}, nil
}
