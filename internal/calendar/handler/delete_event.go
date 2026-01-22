package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func (h *CalendarHttpHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := uuid.Parse(chi.URLParam(r, "eventId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// UserId
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// service
	err = h.eventCommand.Delete(r.Context(), command.EventDeleteInput{
		EventId: eventId,
		UserId:  userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
	}
	render.NoContent(w, r)
}
