package command

import "github.com/TokujouKaisenDonburi/optical-backend/internal/todo/repository"

type TodoCommand struct {
	todoRepository repository.TodoRepository
}

func NewTodoCommand(todoRepository repository.TodoRepository) *TodoCommand {
	if todoRepository == nil {
		panic("todoRepository is nil")
	}
	return &TodoCommand{
		todoRepository: todoRepository,
	}
}
