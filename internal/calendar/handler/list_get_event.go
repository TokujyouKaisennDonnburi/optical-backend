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

type ListGetEventResponse struct {
	Id         string `json:"id"`
	CalendarId string `json:"calendar_id"`
	Title      string `json:"title"`
	Memo       string `json:"memo"`
	Color      string `json:"color"`
	IsAllDay   bool   `json:"is_all_day"`
	StartAt    string `json:"start_at"`
	EndAt      string `json:"end_at"`
	CreatedAt  string `json:"created_at"`
}

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
	response := make([]ListGetEventResponse, 0, len(events))
	for _, event := range events {
		response = append(response, ListGetEventResponse{
			Id:         event.Id.String(),
			CalendarId: event.CalendarId.String(),
			Title:      event.Title,
			Memo:       event.Memo,
			Color:      event.Color,
			IsAllDay:   event.IsAllDay,
			StartAt:    event.StartAt,
			EndAt:      event.EndAt,
			CreatedAt:  event.CreatedAt,
		})
	}
	render.JSON(w, r, response)
}
