package repository

import (
	"context"
	"mime/multipart"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
)

type AvatarRepository interface {
	// Imageを保存する
	Save(ctx context.Context, avatar *user.Avatar) error
	// 画像をアップロードし、URLを返す
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
}
