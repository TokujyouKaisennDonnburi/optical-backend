package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type SchedulerEventRequest struct {
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	IsAllDay  bool      `json:"is_allday"`
	IsDone    bool      `json:"is_done"`
}

func (h *SchedulerHttpHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var request SchedulerEventRequest
	// body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// calendarId
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// schedulerId
	schedulerId, err := uuid.Parse(chi.URLParam(r, "schedulerId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// userId
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrUnauthorized(err))
		return
	}
	err = h.schedulerCommand.SchedulerEvent(r.Context(), command.SchedulerEventInput{
		CalendarId:  calendarId,
		SchedulerId: schedulerId,
		UserId:      userId,
		Date:        request.Date,
		StartTime:   request.StartTime,
		EndTime:     request.EndTime,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	render.NoContent(w, r)
}
