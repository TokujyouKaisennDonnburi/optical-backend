package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type MemberCreateRequest struct {
	Email      string    `JSON:"email"`
}


func (h *CalendarHttpHandler) CreateMembers(w http.ResponseWriter, r *http.Request) {
	// UserId
	userId, err := handler.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w,r,apperr.ErrInternalServerError(err))
		return
	}
	// CalendarId
	calendarId, err := uuid.Parse(chi.URLParam(r,"calendarId"))
	if err != nil {
		_ = render.Render(w,r,apperr.ErrInternalServerError(err))
		return
	}
	// Email
	var request MemberCreateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w,r,apperr.ErrInvalidRequest(err))
		return
	}
	email, err := user.NewEmail(request.Email)
	if err != nil {
		_ = render.Render(w,r,apperr.ErrInvalidRequest(err))
		return
	}
	// 情報をinput
	err = h.calendarCommand.CreateMember(r.Context(), command.MemberCreateInput{
		UserId:     userId,
		CalendarId: calendarId,
		Email: 		string(email),
	})
	if err != nil {
		_ = render.Render(w,r,apperr.ErrInternalServerError(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
