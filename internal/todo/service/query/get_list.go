package query

import (
	"context"

	"github.com/google/uuid"
)

type TodoListQueryInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}

type TodoListQueryOutput struct {
	Id         uuid.UUID
	UserId     uuid.UUID
	CalendarId uuid.UUID
	Name       string
	Items      []TodoListQueryOutputItem
}

type TodoListQueryOutputItem struct {
	Id     uuid.UUID
	UserId uuid.UUID
	Name   string
	IsDone bool
}

func (q *TodoQuery) GetList(
	ctx context.Context,
	input TodoListQueryInput,
) ([]TodoListQueryOutput, error) {
	todoLists, err := q.todoRepository.FindByCalendarId(ctx, input.UserId, input.CalendarId)
	if err != nil {
		return nil, err
	}
	outputList := make([]TodoListQueryOutput, len(todoLists))
	for i, todoList := range todoLists {
		items := make([]TodoListQueryOutputItem, len(todoList.Items))
		for j, item := range todoList.Items {
			items[j] = TodoListQueryOutputItem{
				Id:     item.Id,
				UserId: item.UserId,
				Name:   item.Name,
				IsDone: item.IsDone,
			}
		}
		outputList[i] = TodoListQueryOutput{
			Id:         todoList.Id,
			UserId:     todoList.UserId,
			CalendarId: todoList.CalendarId,
			Name:       todoList.Name,
			Items:      items,
		}
	}
	return outputList, nil
}
