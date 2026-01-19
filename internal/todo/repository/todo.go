package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo/service/query/output"
	"github.com/google/uuid"
)

type TodoRepository interface {
	FindByCalendarId(
		ctx context.Context,
		userId, calendarId uuid.UUID,
	) ([]output.TodoListQueryOutput, error)
	CreateList(
		ctx context.Context,
		list *todo.List,
	) error
	AddItem(
		ctx context.Context,
		listId uuid.UUID,
		todoItem *todo.Item,
	) error
	UpdateItem(
		ctx context.Context,
		userId, itemId uuid.UUID,
		updateFn func(*todo.Item) (*todo.Item, error),
	) error
}
