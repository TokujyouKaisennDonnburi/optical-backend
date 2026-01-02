package handler

import "github.com/TokujouKaisenDonburi/optical-backend/internal/notice/service/query"

// QueryのDI
type NoticeHttpHandler struct {
	NoticeQuery *query.NoticeQuery
}

func NewNoticeHttpHandler(
	noticeQuery *query.NoticeQuery,
) *NoticeHttpHandler {
	if noticeQuery == nil {
		panic("NoticeQuery is nil")
	}
	return &NoticeHttpHandler{
		NoticeQuery: noticeQuery,
	}
}
