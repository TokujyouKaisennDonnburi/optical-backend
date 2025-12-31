package creator

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice/repository"
	"github.com/google/uuid"
)

type NoticeCreator struct {
	repo repository.NoticeRepository
}

func NewNoticeCreator(repo repository.NoticeRepository) *NoticeCreator {
	return &NoticeCreator{
		repo: repo,
	}
}

// オプション型
type CreateOption func(*createOptions)

type createOptions struct {
	eventID    uuid.NullUUID
	calendarID uuid.NullUUID
}

func WithEventID(id uuid.UUID) CreateOption {
	return func(o *createOptions) {
		o.eventID = uuid.NullUUID{UUID: id, Valid: true}
	}
}

func WithCalendarID(id uuid.UUID) CreateOption {
	return func(o *createOptions) {
		o.calendarID = uuid.NullUUID{UUID: id, Valid: true}
	}
}

// 通知作成
// オプションでcreator.WithEventID(eventID), creator.WithCalendarID(calendarID)を指定可能
func (c *NoticeCreator) CreateNotice(
	ctx context.Context,
	userID uuid.UUID,
	title, content string,
	options ...CreateOption,
) error {
	o := &createOptions{}
	for _, opt := range options {
		opt(o)
	}

	n, err := notice.NewNotice(userID, o.eventID, o.calendarID, title, content)
	if err != nil {
		return err
	}

	return c.repo.CreateNotice(ctx, n)
}
