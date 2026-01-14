package handler

import (
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)

type EventTodayResponse struct {
	Date  string                   `json:"date"`
	Items []EventTodayResponseItem `json:"items"`
}

type EventTodayResponseItem struct {
	CalendarId    string    `json:"calendarId"`
	CalendarName  string    `json:"calendarName"`
	CalendarColor string    `json:"calendarColor"`
	Id            string    `json:"id"`
	Title         string    `json:"title"`
	Location      string    `json:"location"`
	Memo          string    `json:"memo"`
	StartAt       time.Time `json:"startAt"`
	EndAt         time.Time `json:"endAt"`
	IsAllDay      bool      `json:"isAllDay"`
}

func (h *CalendarHttpHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	output, err := h.eventQuery.GetToday(r.Context(), query.EventTodayQueryInput{
		UserId: userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	items := make([]EventTodayResponseItem, len(output.Items))
	for i, item := range output.Items {
		items[i] = EventTodayResponseItem{
			CalendarId:    item.CalendarId.String(),
			CalendarName:  item.CalendarName,
			CalendarColor: item.CalendarColor,
			Id:            item.Id.String(),
			Title:         item.Title,
			Location:      item.Location,
			Memo:          item.Memo,
			StartAt:       item.StartAt,
			EndAt:         item.EndAt,
			IsAllDay:      item.IsAllDay,
		}
	}
	render.JSON(w, r, &EventTodayResponse{
		Date:  output.Date,
		Items: items,
	})
}
