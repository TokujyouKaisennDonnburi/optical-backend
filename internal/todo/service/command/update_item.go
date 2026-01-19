package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/google/uuid"
)

type TodoUpdateItemInput struct {
	ItemId uuid.UUID
	UserId uuid.UUID
	Name   string
	IsDone bool
}

func (c *TodoCommand) UpdateItem(
	ctx context.Context,
	input TodoUpdateItemInput,
) error {
	return c.todoRepository.UpdateItem(ctx, input.UserId, input.ItemId, func(item *todo.Item) (*todo.Item, error) {
		item.SetDone(input.IsDone)
		err := item.SetName(input.Name)
		if err != nil {
			return nil, err
		}
		return item, nil
	})
}
