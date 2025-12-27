package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice/service/query/output"
	"github.com/google/uuid"
)

// dbタグでマッピングするためのモデル
type NoticeListQueryModel struct {
	Id         uuid.UUID     `db:"id"`
	UserId     uuid.UUID     `db:"user_id"`
	EventId    uuid.NullUUID `db:"event_id"`
	CalendarId uuid.NullUUID `db:"calendar_id"`
	Title      string        `db:"title"`
	Content    string        `db:"content"`
	IsRead     bool          `db:"is_read"`
	CreatedAt  string        `db:"created_at"`
}

func (r NoticePsqlRepository) ListNoticesByUserId(
	ctx context.Context,
	userId uuid.UUID,
) ([]output.NoticeQueryOutput, error) {
	query := `
		SELECT id, user_id, event_id, calendar_id, title, content, is_read, created_at
		FROM notice
		WHERE user_id = $1
		ORDER BY id DESC
		`

	// クエリ実行
	var rows []NoticeListQueryModel
	err := r.db.SelectContext(ctx, &rows, query, userId)
	if err != nil {
		return nil, err
	}

	// 出力形式に変換
	notices := make([]output.NoticeQueryOutput, len(rows))
	for i, row := range rows {

		notices[i] = output.NoticeQueryOutput{
			Id:         row.Id,
			UserId:     row.UserId,
			EventId:    row.EventId,
			CalendarId: row.CalendarId,
			Title:      row.Title,
			Content:    row.Content,
			IsRead:     row.IsRead,
			CreatedAt:  row.CreatedAt,
		}
	}

	return notices, nil
}
