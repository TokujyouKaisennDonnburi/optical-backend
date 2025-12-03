package handler

// GET /calendars/{calendarId}/events

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// GET /calendars/{calendarId}/events
func (h *CalendarHttpHandler) ListGetEvents(w http.ResponseWriter, r *http.Request) {
	// 1. JWTからuserIDを取得
	userId, err := handler.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	// 2. パスパラメータからcalendarIDを取得
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	// 3. Queryを実行
	events, err := h.eventQuery.ListGetEvents(r.Context(), query.EventQueryInput{
		CalendarID: calendarId,
		UserID:     userId,
	})
	if err != nil {
		// 権限エラーの場合は403を返す
		if err == query.ErrCalendarNotBelongToUser {
			_ = render.Render(w, r, apperr.ErrForbidden(err))
			return
		}
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	// 4. レスポンスを返す
	render.JSON(w, r, events)
}
