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

type SchedulerRequest struct {
	Id           uuid.UUID             `json:"id"`
	CalendarId   uuid.UUID             `json:"calendar_id"`
	UserId       uuid.UUID             `json:"user_id"`
	Title        string                `json:"title"`
	Memo         string                `json:"memo"`
	LimitTime    time.Time             `json:"limit_time"`
	IsAllDay     bool                  `json:"is_all_day"`
	PossibleDate []PossibleDateRequest `json:"possible_date"`
}
type PossibleDateRequest struct {
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (h *SchedulerHttpHandler) GetScheduler(w http.ResponseWriter, r *http.Request) {
	// userId
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrUnauthorized(err))
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
	result, err := h.schedulerQuery.SchedulerQuery(r.Context(), query.SchedulerQueryInput{
		SchedulerId: schedulerId,
		UserId:      userId,
		CalendarId:  calendarId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	responseDates := make([]PossibleDateRequest, len(result.PossibleDate))
	for i, v := range result.PossibleDate {
		responseDates[i] = PossibleDateRequest{
			Date:      v.Date,
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
		}
	}
	responseResult := SchedulerRequest{
		Id:           result.Id,
		CalendarId:   result.CalendarId,
		UserId:       result.UserId,
		Title:        result.Title,
		Memo:         result.Memo,
		LimitTime:    result.LimitTime,
		IsAllDay:     result.IsAllDay,
		PossibleDate: responseDates,
	}
	// response
	render.JSON(w, r, responseResult)
}
