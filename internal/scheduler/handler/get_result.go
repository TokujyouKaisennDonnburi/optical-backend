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

type SchedulerResultResponse struct {
	OwnerId   uuid.UUID        `json:"owner_id"`
	Title     string           `json:"title"`
	Memo      string           `json:"memo"`
	LimitTime time.Time        `json:"limit_time"`
	IsAllDay  bool             `json:"is_allday"`
	Members   []MemberResponse `json:"members"`
	Date      []DateResponse   `json:"date"`
}
type MemberResponse struct {
	UserId   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
}
type DateResponse struct {
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (h *SchedulerHttpHandler) SchedulerResultHandler(w http.ResponseWriter, r *http.Request) {
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
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	result, err := h.schedulerQuery.SchedulerResult(r.Context(), query.SchedulerResultInput{
		CalendarId:  calendarId,
		SchedulerId: schedulerId,
		UserId:      userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	// response
	render.JSON(w, r, result)
}
