package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
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

func (h *CalendarHttpHandler) UpdateCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	var request CalendarUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	// 更新前のカレンダー情報を取得（差分検出用）
	oldCalendar, err := h.calendarQuery.GetCalendar(ctx, query.GetCalendarInput{
		UserId:     userId,
		CalendarId: calendarId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}

	// カレンダー更新
	err = h.calendarCommand.UpdateCalendar(ctx, command.CalendarUpdateInput{
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

	// カレンダー名が変更された場合、メンバーへ通知
	if oldCalendar.Name != request.Name {
		_ = h.calendarNoticeService.NotifyCalendarNameUpdated(ctx, calendarId, oldCalendar.Name, request.Name, userId)
	}

	// カラーが変更された場合、メンバーへ通知
	if string(oldCalendar.Color) != request.Color {
		_ = h.calendarNoticeService.NotifyCalendarColorUpdated(ctx, calendarId, oldCalendar.Name, string(oldCalendar.Color), request.Color, userId)
	}

	render.NoContent(w, r)
}
