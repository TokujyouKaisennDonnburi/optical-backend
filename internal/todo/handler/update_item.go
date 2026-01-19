package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type UpdateTodoItemRequest struct {
	Name   string `json:"name"`
	IsDone bool   `json:"isDone"`
}

func (h *TodoHttpHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	todoItemId, err := uuid.Parse(chi.URLParam(r, "todoItemId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	var request UpdateTodoItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	err = h.todoCommand.UpdateItem(r.Context(), command.TodoUpdateItemInput{
		UserId: userId,
		ItemId: todoItemId,
		Name:   request.Name,
		IsDone: request.IsDone,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.NoContent(w, r)
}
