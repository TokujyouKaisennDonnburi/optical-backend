package output

import "github.com/google/uuid"

type TodoListQueryOutput struct {
	Id            uuid.UUID
	UserId        uuid.UUID
	UserAvatarUrl string
	CalendarId    uuid.UUID
	Name          string
	Items         []TodoListQueryOutputItem
}

type TodoListQueryOutputItem struct {
	Id            uuid.UUID
	UserId        uuid.UUID
	UserAvatarUrl string
	Name          string
	IsDone        bool
}
