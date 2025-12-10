package handler

import (
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query/output"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)

type EventMonthResponse struct {
	Date  string                   `json:"date"`
	Items []EventMonthResponseItem `json:"items"`
}

type EventMonthResponseItem struct {
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

func (h *CalendarHttpHandler) GetByMonth(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		err = render.Render(w, r, apperr.ErrInternalServerError(err))
		if err != nil {
			_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		}
		return
	}
	var output *output.EventMonthQueryOutput
	month := r.URL.Query().Get("month")
	if month == "" {
		output, err = h.eventQuery.GetThisMonth(r.Context(), query.EventThisMonthQueryInput{
			UserId: userId,
		})
	} else {
		output, err = h.eventQuery.GetByMonth(r.Context(), query.EventMonthQueryInput{
			UserId: userId,
			Month:  month,
		})
	}
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	items := make([]EventMonthResponseItem, len(output.Items))
	for i, item := range output.Items {
		items[i] = EventMonthResponseItem{
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
	render.JSON(w, r, &EventMonthResponse{
		Date:  output.Date,
		Items: items,
	})
}
