package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/render"
)

type CalendarListResponse struct {
	Calendars []CalendarResponse `json:"calendars"`
}

type CalendarResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// カレンダー一覧を取得する
func (h *CalendarHttpHandler) GetCalendars(w http.ResponseWriter, r *http.Request) {
	// ユーザーIDを取得
	userId, err := handler.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// カレンダー一覧を取得
	output, err := h.calendarQuery.GetCalendars(r.Context(), query.CalendarQueryInput{
		UserId: userId,
	})
	if err != nil {
		apperr.HandleAppError(w,r,err)
		return
	}
	// レスポンスに変換
	calendars := make([]CalendarResponse, len(output))
	for i, cal := range output {
		calendars[i] = CalendarResponse{
			Id:    cal.Id.String(),
			Name:  cal.Name,
			Color: cal.Color,
		}
	}
	render.JSON(w, r, calendars)
}
