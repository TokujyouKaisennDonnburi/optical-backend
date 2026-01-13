package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)

type SearchEventResponse struct {
	Items []SearchEventResponseItem `json:"items"`
	Total int                       `json:"total"`
	Limit int                       `json:"limit"`
}

type SearchEventResponseItem struct {
	CalendarId    string    `json:"calendarId"`
	CalendarName  string    `json:"calendarName"`
	CalendarColor string    `json:"calendarColor"`
	EventId       string    `json:"eventId"`
	EventTitle    string    `json:"eventTitle"`
	Location      string    `json:"location"`
	Memo          string    `json:"memo"`
	StartAt       time.Time `json:"startAt"`
	EndAt         time.Time `json:"endAt"`
	IsAllDay      bool      `json:"isAllDay"`
}

// イベント検索
func (h *CalendarHttpHandler) SearchEvents(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	// クエリパラメータ取得
	q := r.URL.Query()
	searchQuery := q.Get("query")
	if searchQuery == "" {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(fmt.Errorf("query parameter is required")))
		return
	}

	// 日時パラメータのパース
	var startFrom, startTo time.Time
	if v := q.Get("start_from"); v != "" {
		startFrom, _ = time.Parse(time.RFC3339, v)
	}
	if v := q.Get("start_to"); v != "" {
		startTo, _ = time.Parse(time.RFC3339, v)
	}

	// ページネーションパラメータ
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	output, err := h.eventQuery.SearchEvents(r.Context(), query.SearchEventQueryInput{
		UserId:    userId,
		Query:     searchQuery,
		StartFrom: startFrom,
		StartTo:   startTo,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}

	items := make([]SearchEventResponseItem, len(output.Items))
	for i, item := range output.Items {
		items[i] = SearchEventResponseItem{
			CalendarId:    item.CalendarId.String(),
			CalendarName:  item.CalendarName,
			CalendarColor: item.CalendarColor,
			EventId:       item.EventId.String(),
			EventTitle:    item.EventTitle,
			Location:      item.Location,
			Memo:          item.Memo,
			StartAt:       item.StartAt,
			EndAt:         item.EndAt,
			IsAllDay:      item.IsAllDay,
		}
	}

	render.JSON(w, r, &SearchEventResponse{
		Items: items,
		Total: output.Total,
		Limit: output.Limit,
	})
}
