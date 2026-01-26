package handler

import (
	"net/http"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// 返却
type GetMemberResponse struct {
	UserId   uuid.UUID `json:"userId"`
	Name     string    `json:"name"`
	JoinedAt time.Time `json:"joinedAt,omitempty"`
}

// メンバー一覧
func (h *CalendarHttpHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
	// user認証
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	// クエリパラメータからcalendarIDを取得
	calendarId, err := uuid.Parse(chi.URLParam(r, "calendarId"))
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInvalidRequest(err))
		return
	}

	// メンバーの取得
	outputs, err := h.memberQuery.GetMembers(r.Context(), query.MemberQueryInput{
		UserId:     userId,
		CalendarId: calendarId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}

	// レスポンスデータの作成
	items := make([]GetMemberResponse, len(outputs))
	for i, o := range outputs {
		items[i] = GetMemberResponse{
			UserId:   o.UserId,
			Name:     o.Name,
			JoinedAt: o.JoinedAt,
		}
	}

	render.JSON(w, r, items)
}
