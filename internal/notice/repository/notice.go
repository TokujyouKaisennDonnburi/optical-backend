package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice/service/query/output"
	"github.com/google/uuid"
)

type NoticeRepository interface {

	// 一覧取得
	ListNoticesByUserId(
		ctx context.Context,
		userId uuid.UUID,
	) ([]output.NoticeQueryOutput, error)
}
