package repository

import (
	"context"
	"mime/multipart"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
)

type ImageRepository interface {
	// Imageを保存する
	Save(ctx context.Context, image *calendar.Image) error
	// 画像をアップロードし、URLを返す
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
}
