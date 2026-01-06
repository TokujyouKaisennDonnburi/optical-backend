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

type SchedulerCreateRequest struct {
	Title     string    `json:"title"`
	Memo      string    `json:"memo"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	LimitTime time.Time `json:"limitTime"`
	IsAllDay  bool      `json:"isAllDay"`
}
type SchedulerCreateResponse struct {
	Id uuid.UUID `json:"schedulerId"`
}

func (h *SchedulerHttpHandler) SchedulerCreate(w http.ResponseWriter, r *http.Request) {
	var request SchedulerCreateRequest
	// body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// calendarId
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// userId
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// service
	result, err := h.schedulerCommand.CreateScheduler(r.Context(), command.SchedulerCreateInput{
		CalendarId: calendarId,
		UserId:     userId,
		Title:      request.Title,
		Memo:       request.Memo,
		StartTime:  request.StartTime,
		EndTime:    request.EndTime,
		LimitTime:  request.LimitTime,
		IsAllDay:   request.IsAllDay,
	})
	if err != nil {
		_ = render.Render(w, r, apperr.ErrUnauthorized(err))
		return
	}
	// response
	response := SchedulerCreateResponse{
		Id: result.Id,
	}
	render.JSON(w, r, response)
}
