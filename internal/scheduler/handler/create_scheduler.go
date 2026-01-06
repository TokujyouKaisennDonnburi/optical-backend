package handler

import (
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
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
	Id         uuid.UUID `json:"id"`
	CalendarId uuid.UUID `json:"calendarId"`
	Title      string    `json:"title"`
	Memo       string    `json:"memo"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	IsAllDay   bool      `json:"isAllDay"`
}

func (h *SchedulerHttpHandler) SchedulerCreate(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
}
