package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type AddAttendanceRequest struct {
	SchedulerId uuid.UUID       `json:"scheduler_id"`
	UserId      uuid.UUID       `json:"user_id"`
	Comment     string          `json:"comment"`
	Status      []StatusRequest `json:"status"`
}

type StatusRequest struct {
	Date   time.Time `json:"date"`
	Status int8      `json:"status"`
}

func (h *SchedulerHttpHandler) AddAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request AddAttendanceRequest
	// body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// calendarId
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendar_id"))
	if err != nil {
		render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// schedulerId
	schedulerId, err := uuid.Parse(chi.URLParam(r, "scheduler_id"))
	if err != nil {
		render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// array bind
	Status := make([]command.StatusInput, len(request.Status))
	for i, v := range request.Status {
		Status[i] = command.StatusInput{
			Date:   v.Date,
			Status: v.Status,
		}
	}
	// service
	err = h.schedulerCommand.AddAttendanceCommand(ctx, command.AttendanceInput{
		CalendarId:  calendarId,
		SchedulerId: schedulerId,
		UserId:      request.UserId,
		Comment:     request.Comment,
		Status:      Status,
	})
	if err != nil {
		render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	render.JSON(w, r, nil)
}
