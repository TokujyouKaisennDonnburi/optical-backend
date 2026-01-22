package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CalendarEventInput struct {
	EventId uuid.UUID `json:"event_id"`
}

func (h *CalendarHttpHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	// UserId
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// body
	var request CalendarEventInput
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// service
	err = h.eventCommand.Delete(r.Context(), command.EventDeleteInput{
		EventId: request.EventId,
		UserId:  userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
	}
	render.NoContent(w, r)
}
