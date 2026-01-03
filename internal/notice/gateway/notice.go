package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice"
	"github.com/jmoiron/sqlx"
)

type NoticePsqlRepository struct {
	db *sqlx.DB
}

func NewNoticePsqlRepository(db *sqlx.DB) *NoticePsqlRepository {
	return &NoticePsqlRepository{
		db: db,
	}
}

// 通知作成
func (r NoticePsqlRepository) CreateNotice(
	ctx context.Context,
	notice *notice.Notice,
) error {
	query := `
		INSERT INTO notice(id, user_id, event_id, calendar_id, title, content, is_read, created_at)
		VALUES(:id, :userId, :eventId, :calendarId, :title, :content, :isRead, :createdAt)
	`
	_, err := r.db.NamedExecContext(ctx, query, map[string]any{
		"id":         notice.Id,
		"userId":     notice.UserId,
		"eventId":    notice.EventId,
		"calendarId": notice.CalendarId,
		"title":      notice.Title,
		"content":    notice.Content,
		"isRead":     notice.IsRead,
		"createdAt":  time.Now(),
	})
	return err
}
