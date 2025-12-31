package repository

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/notice"
	"github.com/google/uuid"
)

type NoticeRepository interface {

	// 一覧取得
	ListNoticesByUserId(
		ctx context.Context,
		userId uuid.UUID,
	) ([]notice.Notice, error)

	// 作成
	CreateNotice(
		ctx context.Context,
		notice *notice.Notice,
	) error
}
