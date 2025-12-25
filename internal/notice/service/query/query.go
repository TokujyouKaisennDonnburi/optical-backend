package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice/service/query/output"
	"github.com/google/uuid"
)

type NoticeQuery struct {
	noticeRepository repository.NoticeRepository
}

func NewNoticeQuery(noticeRepo repository.NoticeRepository) *NoticeQuery {
	return &NoticeQuery{
		noticeRepository: noticeRepo,
	}
}

// Input: 一覧取得に必要なデータ
type NoticeListQueryInput struct {
	UserID uuid.UUID
}

// 通知一覧取得
func (q *NoticeQuery) ListGetNotices(ctx context.Context, input NoticeListQueryInput) ([]output.NoticeQueryOutput, error) {

	notices, err := q.noticeRepository.ListNoticesByUserId(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return notices, nil
}
