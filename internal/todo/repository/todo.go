package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/google/uuid"
)

type TodoRepository interface {
	FindByCalendarId(
		ctx context.Context,
		userId, calendarId uuid.UUID,
	) ([]todo.List, error)
	CreateList(
		ctx context.Context,
		list *todo.List,
	) error
	AddItem(
		ctx context.Context,
		listId uuid.UUID,
		todoItem *todo.Item,
	) error
}
