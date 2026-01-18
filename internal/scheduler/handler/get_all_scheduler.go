package handler

import (
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type AllSchedulerResponse struct {
	Id         uuid.UUID `json:"scheduler_id"`
	CalendarId uuid.UUID `json:"calendar_id"`
	UserId     uuid.UUID `json:"user_id"`
	Title      string    `json:"title"`
	Memo       string    `json:"memo"`
	LimitTime  time.Time `json:"limit_time"`
	IsAllDay   bool      `json:"is_allday"`
	IsDone     bool      `json:"is_done"`
}

func (h *SchedulerHttpHandler) GetAllScheduler(w http.ResponseWriter, r *http.Request) {
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// userId
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	output, err := h.schedulerQuery.SchedulerQuery(r.Context(), query.SchedulerInput{
		CalendarId: calendarId,
		UserId:     userId,
	})
	response := AllSchedulerResponse{
		Id:         output.Id,
		CalendarId: output.CalendarId,
		UserId:     output.UserId,
		Title:      output.Title,
		Memo:       output.Memo,
		LimitTime:  output.LimitTime,
		IsAllDay:   output.IsAllDay,
		IsDone:     output.IsDone,
	}
	render.JSON(w, r, response)
}
