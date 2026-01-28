package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/google/uuid"
)

type TodoUpdateListInput struct {
	ListId uuid.UUID
	UserId uuid.UUID
	Name   string
}

func (c *TodoCommand) UpdateList(
	ctx context.Context,
	input TodoUpdateListInput,
) error {
	return c.todoRepository.UpdateList(ctx, input.UserId, input.ListId, func(item *todo.List) (*todo.List, error) {
		err := item.SetName(input.Name)
		if err != nil {
			return nil, err
		}
		return item, nil
	})
}
