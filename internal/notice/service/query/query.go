package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice/repository"
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

// 通知取得で渡す出力データ
type NoticeQueryOutput struct {
	Id         uuid.UUID
	UserId     uuid.UUID
	EventId    uuid.NullUUID
	CalendarId uuid.NullUUID
	Title      string
	Content    string
	IsRead     bool
	CreatedAt  string
}

// 通知一覧取得
func (q *NoticeQuery) ListGetNotices(ctx context.Context, input NoticeListQueryInput) ([]NoticeQueryOutput, error) {

	// repositoryからnotice.Noticeを受け取る
	notices, err := q.noticeRepository.ListNoticesByUserId(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	// NoticeQueryOutputにマッピング
	output := make([]NoticeQueryOutput, len(notices))
	for i, n := range notices {
		output[i] = NoticeQueryOutput{
			Id:         n.Id,
			UserId:     n.UserId,
			EventId:    n.EventId,
			CalendarId: n.CalendarId,
			Title:      n.Title,
			Content:    n.Content,
			IsRead:     n.IsRead,
			CreatedAt:  n.CreatedAt,
		}
	}

	return output, nil
}
