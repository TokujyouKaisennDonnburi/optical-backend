package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/google/uuid"
)

type TodoCreateInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
	Name       string
}

type TodoCreateOutput struct {
	Id         uuid.UUID
	UserId     uuid.UUID
	CalendarId uuid.UUID
	Name       string
	Items      []TodoCreateOutputItem
}

type TodoCreateOutputItem struct {
	Id     uuid.UUID
	UserId uuid.UUID
	Name   string
	IsDone bool
}

// TODOリストを新規作成する
func (c *TodoCommand) CreateList(ctx context.Context, input TodoCreateInput) (*TodoCreateOutput, error) {
	todoList, err := todo.NewList(input.UserId, input.CalendarId, input.Name)
	if err != nil {
		return nil, err
	}
	err = c.todoRepository.CreateList(ctx, todoList)
	if err != nil {
		return nil, err
	}
	return &TodoCreateOutput{
		Id:         todoList.Id,
		UserId:     todoList.UserId,
		CalendarId: todoList.CalendarId,
		Name:       todoList.Name,
		Items:      []TodoCreateOutputItem{},
	}, nil
}
