package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type MemberCreateRequest struct {
	Email []string `json:"email"`
}

func (h *CalendarHttpHandler) CreateMembers(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	var request MemberCreateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	err = h.calendarCommand.InviteMember(r.Context(), command.MemberCreateInput{
		UserId:     userId,
		CalendarId: calendarId,
		Emails:     request.Email,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
