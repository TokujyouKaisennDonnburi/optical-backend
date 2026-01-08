package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
)

type TodoRepository interface {
	CreateList(
		ctx context.Context,
		list *todo.List,
	) error
}
