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
	LimitTime time.Time `json:"limitTime"`
	IsAllDay  bool      `json:"isAllDay"`
	Dates     []SchedulerCreateDateRequest `json:"dates"`
}

type SchedulerCreateDateRequest struct {
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
type SchedulerCreateResponse struct {
	Id uuid.UUID `json:"schedulerId"`
}

func (h *SchedulerHttpHandler) CreateScheduler(w http.ResponseWriter, r *http.Request) {
	var request SchedulerCreateRequest
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
	// userId
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrUnauthorized(err))
		return
	}
	// array
	dates := make([]command.SchedulerCreateDateInput, len(request.Dates))
	for i, date := range request.Dates {
		dates[i] = command.SchedulerCreateDateInput{
			Date:      date.Date,
			StartTime: date.StartTime,
			EndTime:   date.EndTime,
		}
	}
	// service
	result, err := h.schedulerCommand.CreateScheduler(r.Context(), command.SchedulerCreateInput{
		CalendarId: calendarId,
		UserId:     userId,
		Title:      request.Title,
		Memo:       request.Memo,
		LimitTime:  request.LimitTime,
		IsAllDay:   request.IsAllDay,
		Dates:      dates,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	// response
	response := SchedulerCreateResponse{
		Id: result.Id,
	}
	render.JSON(w, r, response)
}
