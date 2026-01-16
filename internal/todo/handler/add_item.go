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

type CreateTodoItemRequest struct {
	Name string `json:"name"`
}

type CreateTodoItemResponse struct {
	Id uuid.UUID `json:"id"`
}

func (h *TodoHttpHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	todoListId, err := uuid.Parse(chi.URLParam(r, "todoListId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	var request CreateTodoItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	output, err := h.todoCommand.AddItem(r.Context(), command.TodoCreateItemInput{
		UserId: userId,
		ListId: todoListId,
		Name:   request.Name,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.JSON(w, r, CreateTodoItemResponse{
		Id: output.Id,
	})
}
