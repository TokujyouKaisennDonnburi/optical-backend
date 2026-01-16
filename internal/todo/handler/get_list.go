package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type TodoListResponse struct {
	Id         uuid.UUID              `json:"id"`
	UserId     uuid.UUID              `json:"userId"`
	CalendarId uuid.UUID              `json:"calendarId"`
	Name       string                 `json:"name"`
	Items      []TodoListResponseItem `json:"items"`
}

type TodoListResponseItem struct {
	Id     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"userId"`
	Name   string    `json:"name"`
	IsDone bool      `json:"isDone"`
}

func (h *TodoHttpHandler) GetList(w http.ResponseWriter, r *http.Request) {
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
	outputList, err := h.todoQuery.GetList(r.Context(), query.TodoListQueryInput{
		UserId:     userId,
		CalendarId: calendarId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	responseList := make([]TodoListResponse, len(outputList))
	for i, output := range outputList {
		responseItems := make([]TodoListResponseItem, len(output.Items))
		for j, item := range output.Items {
			responseItems[j] = TodoListResponseItem{
				Id:     item.Id,
				UserId: item.UserId,
				Name:   item.Name,
				IsDone: item.IsDone,
			}
		}
		responseList[i] = TodoListResponse{
			Id:         output.Id,
			UserId:     output.UserId,
			CalendarId: output.CalendarId,
			Name:       output.Name,
			Items:      responseItems,
		}
	}
	render.JSON(w, r, responseList)
}
