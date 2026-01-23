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

type AttendanceResponse struct {
	UserId  uuid.UUID        `json:"userId"`
	Comment string           `json:"comment"`
	Status  []StatusResponse `json:"status"`
}
type StatusResponse struct {
	Date   time.Time `json:"date"`
	Status int8      `json:"status"`
}

func (h *SchedulerHttpHandler) GetAttendance(w http.ResponseWriter, r *http.Request) {
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
	// service
	output, err := h.schedulerQuery.AttendanceQuery(r.Context(), query.AttendanceQueryInput{
		CalendarId:  calendarId,
		SchedulerId: schedulerId,
		UserId:      userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	// array bind
	statuses := make([]StatusResponse, len(output.Status))
	for i, v := range output.Status {
		statuses[i] = StatusResponse{
			Date:   v.Date,
			Status: v.Status,
		}
	}
	// bind
	response := AttendanceResponse{
		UserId:  output.UserId,
		Comment: output.Comment,
		Status:  statuses,
	}
	// response
	render.JSON(w, r, response)
}
