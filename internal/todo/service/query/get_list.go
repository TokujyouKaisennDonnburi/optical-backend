package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo/service/query/output"
	"github.com/google/uuid"
)

type TodoListQueryInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}

func (q *TodoQuery) GetList(
	ctx context.Context,
	input TodoListQueryInput,
) ([]output.TodoListQueryOutput, error) {
	return q.todoRepository.FindByCalendarId(ctx, input.UserId, input.CalendarId)
}
