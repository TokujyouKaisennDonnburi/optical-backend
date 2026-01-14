package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type EventUpdateRequest struct {
	Title     string    `json:"title"`
	Memo      string    `json:"memo"`
	Location  string    `json:"location"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	IsAllDay  bool      `json:"isAllDay"`
}

func (h *CalendarHttpHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	eventId, err := uuid.Parse(chi.URLParam(r, "eventId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	var request EventUpdateRequest
	// リクエストJSONをバインド
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	err = h.eventCommand.UpdateEvent(r.Context(), command.EventUpdateInput{
		UserId:    userId,
		EventId:   eventId,
		Title:     request.Title,
		Memo:      request.Memo,
		Location:  request.Location,
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
		IsAllDay:  request.IsAllDay,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.NoContent(w, r)
}
