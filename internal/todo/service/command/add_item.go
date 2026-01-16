package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/google/uuid"
)

type TodoCreateItemInput struct {
	UserId uuid.UUID
	ListId uuid.UUID
	Name   string
}

type TodoCreateItemOutput struct {
	Id uuid.UUID
}

func (c *TodoCommand) AddItem(
	ctx context.Context,
	input TodoCreateItemInput,
) (*TodoCreateItemOutput, error) {
	todoItem, err := todo.NewItem(input.ListId, input.UserId, input.Name)
	if err != nil {
		return nil, err
	}
	err = c.todoRepository.AddItem(ctx, input.ListId, todoItem)
	if err != nil {
		return nil, err
	}
	return &TodoCreateItemOutput{
		Id: todoItem.Id,
	}, nil
}
