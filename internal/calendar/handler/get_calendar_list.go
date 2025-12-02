package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)


type CalendarListResponse struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Color    string  `json:"color"`
	ImageUrl *string `json:"imageUrl"`
}

// カレンダー一覧を取得する
func (h *CalendarHttpHandler) GetCalendars(w http.ResponseWriter, r *http.Request) {
	// ユーザーIDを取得
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// カレンダー一覧を取得
	output, err := h.calendarQuery.GetCalendars(r.Context(), query.CalendarQueryInput{
		UserId: userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	// レスポンスに変換
	calendars := make([]CalendarListResponse, len(output))
	for i, cal := range output {
		var imageUrl *string = nil
		if cal.Image.Valid {
			imageUrl = &cal.Image.Url
		}
		calendars[i] = CalendarListResponse{
			Id:       cal.Id.String(),
			Name:     cal.Name,
			Color:    cal.Color,
			ImageUrl: imageUrl,
		}
	}
	render.JSON(w, r, calendars)
}
