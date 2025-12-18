package repository

import (
	"context"
	"mime/multipart"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type AvatarRepository interface {
	// Imageを保存する
	Save(ctx context.Context, userId uuid.UUID, avatar *user.Avatar) error
	// 画像をアップロードし、URLを返す
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
}
