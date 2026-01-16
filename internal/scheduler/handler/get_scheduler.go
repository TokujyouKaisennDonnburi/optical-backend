package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func (h *SchedulerHttpHandler) GetScheduler(w http.ResponseWriter, r *http.Request) {
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
	// response
	render.JSON(w, r, result)
}
