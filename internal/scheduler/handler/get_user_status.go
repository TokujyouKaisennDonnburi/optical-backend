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

type UserStatusResponse struct {
	UserId uuid.UUID        `json:"user_id"`
	Status []StatusResponse `json:"status"`
}
type StatusResponse struct {
	Date   time.Time `json:"date"`
	Status int8      `json:"status"`
}

func (h *SchedulerHttpHandler) SchedulerHandler(w http.ResponseWriter, r *http.Request) {
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
	service, err := h.schedulerQuery.UserStatusQuery(r.Context(), query.SchedulerUserStatusInput{
		SchedulerId: schedulerId,
		UserId:      userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
	}
	statuses := make([]StatusResponse, len(service.Status))
	for i, v := range service.Status {
		statuses[i] = StatusResponse{
			Date:   v.Date,
			Status: v.Status,
		}
	}
	response := UserStatusResponse{
		UserId: service.UserId,
		Status: statuses,
	}
	render.JSON(w, r, response)
}
