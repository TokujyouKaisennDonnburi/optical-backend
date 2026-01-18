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

type CreateTodoListRequest struct {
	Name string `json:"name"`
}

type CreateTodoListResponse struct {
	Id         uuid.UUID                    `json:"id"`
	UserId     uuid.UUID                    `json:"userId"`
	CalendarId uuid.UUID                    `json:"calendarId"`
	Name       string                       `json:"name"`
	Items      []CreateTodoListResponseItem `json:"items"`
}

type CreateTodoListResponseItem struct {
	Id     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"userId"`
	Name   string    `json:"name"`
	IsDone bool      `json:"isDone"`
}

func (h *TodoHttpHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	var request CreateTodoListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	output, err := h.todoCommand.CreateList(r.Context(), command.TodoCreateInput{
		UserId:     userId,
		CalendarId: calendarId,
		Name:       request.Name,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	items := make([]CreateTodoListResponseItem, len(output.Items))
	render.JSON(w, r, CreateTodoListResponse{
		Id:         output.Id,
		UserId:     output.UserId,
		CalendarId: output.CalendarId,
		Name:       output.Name,
		Items:      items,
	})
}
