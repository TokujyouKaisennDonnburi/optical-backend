package psql

import (
	"context"
	"database/sql"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TodoListAndItemModel struct {
	Id         uuid.UUID      `db:"id"`
	UserId     uuid.UUID      `db:"user_id"`
	CalendarId uuid.UUID      `db:"calendar_id"`
	Name       string         `db:"name"`
	ItemId     uuid.NullUUID  `db:"item_id"`
	ItemUserId uuid.NullUUID  `db:"item_user_id"`
	ItemName   sql.NullString `db:"item_name"`
	ItemIsDone sql.NullBool   `db:"item_is_done"`
}

func FindTodoListById(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*todo.List, error) {
	query := `
		SELECT lists.id AS id, lists.user_id AS user_id, lists.calendar_id AS calendar_id, lists.name AS name,
			items.id AS item_id, items.user_id AS item_user_id, items.name AS item_name, items.is_done AS item_is_done
		FROM todo_lists lists
		LEFT JOIN todo_items items
			ON lists.id = items.todo_list_id
		WHERE lists.id = $1
	`
	var models []TodoListAndItemModel
	err := tx.SelectContext(ctx, &models, query, id)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, apperr.NotFoundError("todo list not found")
	}
	todoItems := make([]todo.Item, 0)
	for _, model := range models {
		if !model.ItemId.Valid {
			continue
		}
		todoItems = append(todoItems, todo.Item{
			Id:     model.ItemId.UUID,
			UserId: model.ItemUserId.UUID,
			Name:   model.ItemName.String,
			IsDone: model.ItemIsDone.Bool,
		})
	}
	return &todo.List{
		Id:         models[0].Id,
		UserId:     models[0].UserId,
		CalendarId: models[0].CalendarId,
		Name:       models[0].Name,
		Items:      todoItems,
	}, nil
}
