package gateway

import (
	"context"
	"maps"
	"slices"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
)

// TODOリストをカレンダーから取得
func (r *TodoPsqlRepository) FindByCalendarId(
	ctx context.Context,
	userId, calendarId uuid.UUID,
) ([]todo.List, error) {
	query := `
		SELECT
			lists.id, lists.user_id, lists.calendar_id, lists.name
			items.id AS item_id, items.user_id AS item_user_id, items.name AS item_name, items.is_done as item_is_done
		FROM todo_lists lists
		JOIN calendar_members
			ON calendar_members.calendar_id = lists.calendar_id
		LEFT JOIN todo_items items
			ON lists.id = items.todo_list_id
		WHERE 
			calendar_members.user_id = $1
			AND lists.calendar_id = $2
	`
	var todoListModels []psql.TodoListAndItemModel
	err := r.db.SelectContext(ctx, &todoListModels, query, userId, calendarId)
	if err != nil {
		return nil, err
	}
	todoListMap := map[uuid.UUID]todo.List{}
	for _, todoListModel := range todoListModels {
		model, ok := todoListMap[todoListModel.Id]
		if !ok {
			model = todo.List{
				Id:         todoListModel.Id,
				UserId:     todoListModel.UserId,
				CalendarId: todoListModel.CalendarId,
				Name:       todoListModel.Name,
			}
		}
		if !todoListModel.ItemId.Valid {
			continue
		}
		model.Items = append(model.Items, todo.Item{
			Id:     todoListModel.ItemId.UUID,
			ListId: todoListModel.Id,
			UserId: todoListModel.ItemUserId.UUID,
			Name:   todoListModel.ItemName.String,
			IsDone: todoListModel.ItemIsDone.Bool,
		})
		todoListMap[todoListModel.Id] = model
	}
	return slices.Collect(maps.Values(todoListMap)), nil
}
