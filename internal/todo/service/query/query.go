package query

import "github.com/TokujouKaisenDonburi/optical-backend/internal/todo/repository"

type TodoQuery struct {
	todoRepository repository.TodoRepository
}

func NewTodoQuery(todoRepository repository.TodoRepository) *TodoQuery {
	if todoRepository == nil {
		panic("todoRepository is nil")
	}
	return &TodoQuery{
		todoRepository: todoRepository,
	}
}
