package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/handler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type MemberCreateRequest struct {
	UserId     uuid.UUID `JSON:"userId"`
	CalendarId uuid.UUID `JSON:"calendarId"`
	Email      string    `JSON:"email"`
}

type MemberCreateResponse struct {
	UserId 		uuid.UUID
	Name		string
	JoinedAt 	time.Time
}

func (h *CalendarHttpHandler) CreateMembers(w http.ResponseWriter, r *http.Request) {
	// TODO: リクエストから必要な情報を取得する
	//   1. コンテキストからUserId
	//   2. URLパラメータからCalendarId
	//   3. リクエストボディからEmail

	// UserId
	userId, err := handler.GetUserIdFromContext(r)
	if err != nil {
		_ =render.Render(w,r,apperr.ErrInternalServerError(err))
		return
	}
	// CalendarId
	calendarId, err := uuid.Parse(chi.URLParam(r,"calendarId"))
	if err != nil {
		_ =render.Render(w,r,apperr.ErrInternalServerError(err))
		return
	}
	var request MemberCreateRequest
	// Email
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = render.Render(w,r,apperr.ErrInvalidRequest(err))
		return
	}

	// TODO: サービス層を呼び出してメンバーを作成する
	output, err := h.calendarCommand.CreateMember(r.Context(), command.MemberCreateInput{
		UserId:     userId,
		CalendarId: calendarId,
		Email: 		request.Email,
	})
	if err != nil {
		_ = render.Render(w,r,apperr.ErrInternalServerError(err))
		return
	}

	// TODO: レスポンスを作成してJSONで返す
	render.JSON(w, r, MemberCreateResponse{
		UserId:   output.UserId,
		Name:     output.Name,
		JoinedAt: output.JoinedAt,

	})
}
