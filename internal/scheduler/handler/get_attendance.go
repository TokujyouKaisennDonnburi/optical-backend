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

type AttendanceQueryRequest struct {
	SchedulerId uuid.UUID `json:"schedulerId"`
}
type AttendanceQueryRespons struct {
	Id           uuid.UUID              `json:"schedulerId"`
	CalendarId   uuid.UUID              `json:"calendarId"`
	UserId       uuid.UUID              `json:"userId"`
	Title        string                 `json:"title"`
	Memo         string                 `json:"memo"`
	LimitTime    time.Time              `json:"limitTime"`
	IsAllDay     bool                   `json:"isAllDay"`
	PossibleDate []PossibleDateResponse `json:"possibleDate"`
}
type PossibleDateResponse struct {
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

func (h *SchedulerHttpHandler) GetAttendance(w http.ResponseWriter, r *http.Request) {
	// userId
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// calendarId
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// result
	var request AttendanceQueryRequest
	result, err := h.schedulerQuery.AttendanceQuery(r.Context(), query.AttendanceQueryInput{
		SchedulerId: request.SchedulerId,
		UserId:      userId,
		CalendarId:  calendarId,
	})
	// response
	render.JSON(w, r, result)
}
