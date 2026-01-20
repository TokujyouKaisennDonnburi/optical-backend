package gateway

import (
	"context"
	"database/sql"
	"maps"
	"slices"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/storage"
	"github.com/google/uuid"
)

type TodoListAndItemAndAvatarModel struct {
	Id                  uuid.UUID      `db:"id"`
	UserId              uuid.UUID      `db:"user_id"`
	AvatarUrl           sql.NullString `db:"avatar_url"`
	AvatarIsFullURL     sql.NullBool   `db:"avatar_is_full_url"`
	CalendarId          uuid.UUID      `db:"calendar_id"`
	Name                string         `db:"name"`
	ItemId              uuid.NullUUID  `db:"item_id"`
	ItemUserId          uuid.NullUUID  `db:"item_user_id"`
	ItemName            sql.NullString `db:"item_name"`
	ItemAvatarUrl       sql.NullString `db:"item_avatar_url"`
	ItemAvatarIsFullURL sql.NullBool   `db:"item_avatar_is_full_url"`
	ItemIsDone          sql.NullBool   `db:"item_is_done"`
}

// TODOリストをカレンダーから取得
func (r *TodoPsqlRepository) FindByCalendarId(
	ctx context.Context,
	userId, calendarId uuid.UUID,
) ([]output.TodoListQueryOutput, error) {
	query := `
		SELECT
			lists.id, lists.user_id, lists.calendar_id, lists.name, 
			avatars.url AS avatar_url, avatars.is_full_url AS avatar_is_full_url,
			items.id AS item_id, items.user_id AS item_user_id, items.name AS item_name, 
			item_avatars.url AS item_avatar_url, item_avatars.is_full_url AS item_avatar_is_full_url, 
			items.is_done as item_is_done
		FROM todo_lists lists
		JOIN calendar_members
			ON calendar_members.calendar_id = lists.calendar_id
		LEFT JOIN user_profiles
			ON lists.user_id = user_profiles.user_id
		LEFT JOIN avatars
			ON user_profiles.avatar_id = avatars.id
		LEFT JOIN todo_items items
			ON lists.id = items.todo_list_id
		LEFT JOIN user_profiles item_user_profiles
			ON items.user_id = item_user_profiles.user_id
		LEFT JOIN avatars item_avatars
			ON item_user_profiles.avatar_id = item_avatars.id
		WHERE 
			calendar_members.user_id = $1
			AND lists.calendar_id = $2
	`
	var todoListModels []TodoListAndItemAndAvatarModel
	err := r.db.SelectContext(ctx, &todoListModels, query, userId, calendarId)
	if err != nil {
		return nil, err
	}
	todoListMap := map[uuid.UUID]output.TodoListQueryOutput{}
	for _, todoListModel := range todoListModels {
		model, ok := todoListMap[todoListModel.Id]
		userAvatarUrl := todoListModel.AvatarUrl.String
		if todoListModel.AvatarUrl.Valid && todoListModel.AvatarIsFullURL.Bool {
			userAvatarUrl = storage.GetImageStorageBaseUrl() + "/" + userAvatarUrl
		}

		if !ok {
			model = output.TodoListQueryOutput{
				Id:            todoListModel.Id,
				UserId:        todoListModel.UserId,
				UserAvatarUrl: userAvatarUrl,
				CalendarId:    todoListModel.CalendarId,
				Name:          todoListModel.Name,
			}
			todoListMap[todoListModel.Id] = model
		}
		if !todoListModel.ItemId.Valid {
			continue
		}
		itemUserAvatarUrl := todoListModel.ItemAvatarUrl.String
		if todoListModel.ItemAvatarUrl.Valid && todoListModel.ItemAvatarIsFullURL.Bool {
			itemUserAvatarUrl = storage.GetImageStorageBaseUrl() + "/" + itemUserAvatarUrl
		}
		model.Items = append(model.Items, output.TodoListQueryOutputItem{
			Id:            todoListModel.ItemId.UUID,
			UserId:        todoListModel.ItemUserId.UUID,
			Name:          todoListModel.ItemName.String,
			UserAvatarUrl: itemUserAvatarUrl,
			IsDone:        todoListModel.ItemIsDone.Bool,
		})
		todoListMap[todoListModel.Id] = model
	}
	return slices.Collect(maps.Values(todoListMap)), nil
}
