package repository

import (
	"context"
	"mime/multipart"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
)

type ImageRepository interface {
	FindById(ctx context.Context, imageId uuid.UUID) (*calendar.Image, error)
	// Imageを保存する
	Save(ctx context.Context, image *calendar.Image) error
	// 画像をアップロードし、URLを返す
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
}
