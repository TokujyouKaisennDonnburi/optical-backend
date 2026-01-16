package handler

import "github.com/TokujouKaisenDonburi/optical-backend/internal/todo/service/command"

type TodoHttpHandler struct {
	todoCommand *command.TodoCommand
}

func NewTodoHttpHandler(
	todoCommand *command.TodoCommand,
) *TodoHttpHandler {
	if todoCommand == nil {
		panic("todoCommand is nil")
	}
	return &TodoHttpHandler{
		todoCommand: todoCommand,
	}
}
