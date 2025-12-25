package handler

import (
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice/service/query"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/go-chi/render"
)

type NoticeResponse struct {
	Id         string  `json:"id"`
	UserId     string  `json:"userId"`
	EventId    *string `json:"eventId"`
	CalendarId *string `json:"calendarId"`
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	IsRead     bool    `json:"isRead"`
	CreatedAt  string  `json:"createdAt"`
}

func (h *NoticeHttpHandler) GetNotices(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		_ = render.Render(w, r, apperr.ErrInternalServerError(err))
		return
	}

	// 通知一覧を取得
	output, err := h.NoticeQuery.ListGetNotices(r.Context(), query.NoticeListQueryInput{
		UserID: userId,
	})
	if err != nil {
		apperr.HandleAppError(w, r, err)
		return
	}

	// レスポンスに変換
	notice := make([]NoticeResponse, len(output))
	for i, n := range output {
		var eventIdStrPtr, calendarIdStrPtr *string
		if n.EventId != nil {
			s := n.EventId.String()
			eventIdStrPtr = &s
		}
		if n.CalendarId != nil {
			s := n.CalendarId.String()
			calendarIdStrPtr = &s
		}

		notice[i] = NoticeResponse{
			Id:         n.Id.String(),
			UserId:     n.UserId.String(),
			EventId:    eventIdStrPtr,
			CalendarId: calendarIdStrPtr,
			Title:      n.Title,
			Content:    n.Content,
			IsRead:     n.IsRead,
			CreatedAt:  n.CreatedAt,
		}
	}

	render.JSON(w, r, notice)
}
