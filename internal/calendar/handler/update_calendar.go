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
	Name      string   `json:"name"`
	Color     string   `json:"color"`
	ImageId   string   `json:"imageId"`
	Members   []string `json:"members"`
	OptionIds []int32  `json:"optionIds"`
}

// カレンダー更新のためのHTTPハンドラー
func (h *CalendarHttpHandler) UpdateCalendar(w http.ResponseWriter, r *http.Request) {

	// リクエスト
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	// リクエストJSONをGo構造体にバインド
	var request CalendarUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	var ImageId uuid.UUID
	if request.ImageId != "" {
		ImageId, err = uuid.Parse(request.ImageId)
		if err != nil {
			render.Render(w, r, apperr.ErrInvalidRequest(err))
			return
		}
	}

	err = h.calendarCommand.UpdateCalendar(r.Context(), command.CalendarUpdateInput{
		UserId:        userId,
		CalendarId:    calendarId,
		CalendarName:  request.Name,
		CalendarColor: request.Color,
		MemberEmails:  request.Members,
		ImageId:       ImageId,
		OptionIds:     request.OptionIds,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}

	// レスポンス
	render.NoContent(w, r)
}
