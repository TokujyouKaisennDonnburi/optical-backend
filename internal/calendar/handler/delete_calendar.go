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

func (h *CalendarHttpHandler) DeleteCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	// カレンダー削除
	err = h.calendarCommand.DeleteCalendar(ctx, command.CalendarDeleteInput{
		CalendarId: calendarId,
		UserId:     userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}

	render.NoContent(w, r)
}
