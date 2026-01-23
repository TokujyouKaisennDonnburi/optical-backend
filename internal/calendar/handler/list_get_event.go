package handler

// GET /calendars/{calendarId}/events

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ListGetEventResponse struct {
	Id            string `json:"id"`
	CalendarId    string `json:"calendarId"`
	UserId        string `json:"userId"`
	CalendarColor string `json:"calendarColor"`
	Title         string `json:"title"`
	Memo          string `json:"memo"`
	Location      string `json:"location"`
	IsAllDay      bool   `json:"isAllDay"`
	StartAt       string `json:"startAt"`
	EndAt         string `json:"endAt"`
	CreatedAt     string `json:"createdAt"`
}

// GET /calendars/{calendarId}/events
func (h *CalendarHttpHandler) ListGetEvents(w http.ResponseWriter, r *http.Request) {
	// 1. JWTからuserIDを取得
	userId, err := auth.GetUserIdFromContext(r)
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
		apperr.HandleAppError(w, r, err)
		return
	}

	// 4. レスポンスを返す
	response := make([]ListGetEventResponse, 0, len(events))
	for _, event := range events {
		response = append(response, ListGetEventResponse{
			Id:            event.Id.String(),
			CalendarId:    event.CalendarId.String(),
			UserId:        event.UserId.String(),
			CalendarColor: event.CalendarColor,
			Title:         event.Title,
			Memo:          event.Memo,
			Location:      event.Location,
			IsAllDay:      event.IsAllDay,
			StartAt:       event.StartAt,
			EndAt:         event.EndAt,
			CreatedAt:     event.CreatedAt,
		})
	}
	render.JSON(w, r, response)
}
