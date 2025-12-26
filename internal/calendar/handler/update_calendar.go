package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CalendarUpdateRequest struct {
	Name      string  `json:"name"`
	Color     string  `json:"color"`
	OptionIds []int32 `json:"optionIds"`
}

// カレンダー更新のためのHTTPハンドラー
func (h *CalendarHttpHandler) UpdateCalendar(w http.ResponseWriter, r *http.Request) {

	// リクエスト
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	// リクエストJSONをGo構造体にバインド
	var request CalendarUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	// カレンダー更新
	err = h.calendarCommand.UpdateCalendar(r.Context(), command.CalendarUpdateInput{
		UserId:        userId,
		CalendarId:    calendarId,
		CalendarName:  request.Name,
		CalendarColor: request.Color,
		OptionIds:     request.OptionIds,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}

	// レスポンス
	render.NoContent(w, r)
}
