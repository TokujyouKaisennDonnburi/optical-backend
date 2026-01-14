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

type AddAttendanceRequest struct {
	Comment string          `json:"comment"`
	Status  []StatusRequest `json:"status"`
}

type StatusRequest struct {
	Date   time.Time `json:"date"`
	Status int8      `json:"status"`
}

func (h *SchedulerHttpHandler) AddAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	var request AddAttendanceRequest
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
	// array bind
	status := make([]command.StatusInput, len(request.Status))
	for i, v := range request.Status {
		status[i] = command.StatusInput{
			Date:   v.Date,
			Status: v.Status,
		}
	}
	// service
	err = h.schedulerCommand.AddAttendanceCommand(r.Context(), command.AttendanceInput{
		CalendarId:  calendarId,
		SchedulerId: schedulerId,
		UserId:      userId,
		Comment:     request.Comment,
		Status:      status,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
