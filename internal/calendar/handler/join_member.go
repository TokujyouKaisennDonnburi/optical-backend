package handler

import (
	"errors"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// 従来の参加（仮メンバー状態からの参加）
func (h *CalendarHttpHandler) JoinMember(w http.ResponseWriter, r *http.Request) {
	// userId
	user, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// CalendarId
	calendar, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// input
	err = h.calendarCommand.JoinMember(r.Context(), command.CalendarJoinInput{
		UserId:     user,
		CalendarId: calendar,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// トークンを使用した参加 (登録済みユーザーでも未登録ユーザーでもメアドが違っても登録できる)
// 登録完了後に実行
func (h *CalendarHttpHandler) JoinMemberWithToken(w http.ResponseWriter, r *http.Request) {
	// userId
	user, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}
	// CalendarId
	calendar, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}
	// Token
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(errors.New("token is required")))
		return
	}
	token, err := uuid.Parse(tokenStr)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	// input
	err = h.calendarCommand.JoinMemberWithToken(r.Context(), command.JoinWithTokenInput{
		UserId:     user,
		CalendarId: calendar,
		Token:      token,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

