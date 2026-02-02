package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
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
	n *notice.Notice,
) error {
	return db.RunInTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		query := `
			INSERT INTO notice(id, user_id, event_id, calendar_id, title, content, is_read, created_at)
			VALUES(:id, :userId, :eventId, :calendarId, :title, :content, :isRead, :createdAt)
		`
		_, err := tx.NamedExecContext(ctx, query, map[string]any{
			"id":         n.Id,
			"userId":     n.UserId,
			"eventId":    n.EventId,
			"calendarId": n.CalendarId,
			"title":      n.Title,
			"content":    n.Content,
			"isRead":     n.IsRead,
			"createdAt":  time.Now().UTC(),
		})
		return err
	})
}
