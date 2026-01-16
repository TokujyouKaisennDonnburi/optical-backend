package handler

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo/service/query"
)

type TodoHttpHandler struct {
	todoQuery   *query.TodoQuery
	todoCommand *command.TodoCommand
}

func NewTodoHttpHandler(
	todoQuery *query.TodoQuery,
	todoCommand *command.TodoCommand,
) *TodoHttpHandler {
	if todoQuery == nil {
		panic("todoQuery is nil")
	}
	if todoCommand == nil {
		panic("todoCommand is nil")
	}
	return &TodoHttpHandler{
		todoQuery:   todoQuery,
		todoCommand: todoCommand,
	}
}
