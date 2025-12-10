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

type EventCreateRequest struct {
	Title     string    `jsno:"title"`
	Memo      string    `jsno:"memo"`
	Color     string    `jsno:"color"`
	Location  string    `jsno:"location"`
	StartTime time.Time `jsno:"startTime"`
	EndTime   time.Time `jsno:"endTime"`
	IsAllDay  bool      `jsno:"isAllDay"`
}

type EventCreateResponse struct {
	Id string `json:"id"`
}

func (h *CalendarHttpHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
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
	var request EventCreateRequest
	// リクエストJSONをバインド
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	output, err := h.eventCommand.Create(r.Context(), command.EventCreateInput{
		UserId:     userId,
		CalendarId: calendarId,
		Title:      request.Title,
		Memo:       request.Memo,
		Color:      request.Color,
		Location:   request.Location,
		StartTime:  request.StartTime,
		EndTime:    request.EndTime,
		IsAllDay:   request.IsAllDay,
	})
	if err != nil {
		apperr.HandleAppError(w,r,err)
		return
	}
	render.JSON(w, r, EventCreateResponse{
		Id: output.Id.String(),
	})
}
