package gateway

import (
	"context"
	"mime/multipart"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type AvatarPsqlAndMinioRepository struct {
	db          *sqlx.DB
	minioClient *minio.Client
	bucketName  string
}

func NewAvatarPsqlAndMinioRepository(db *sqlx.DB, minioClient *minio.Client, bucketName string) *AvatarPsqlAndMinioRepository {
	if db == nil {
		panic("psql db is nil")
	}
	if minioClient == nil {
		panic("minioClient is nil")
	}
	if bucketName == "" {
		panic("bucketName is empty")
	}
	return &AvatarPsqlAndMinioRepository{
		db:          db,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (r *AvatarPsqlAndMinioRepository) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	filePath := "avatars/images/"
	_, err := r.minioClient.PutObject(ctx, r.bucketName, filePath+header.Filename, file, header.Size, minio.PutObjectOptions{
		ContentType: "image/png",
	})
	if err != nil {
		return "", err
	}
	return filePath + header.Filename, nil
}

func (r *AvatarPsqlAndMinioRepository) Save(ctx context.Context, userId uuid.UUID, avatar *user.Avatar) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		query := `
			INSERT INTO avatars(id, url, is_relative_path)
			VALUES(:id, :url, :isRelativePath)
		`
		_, err := tx.NamedExecContext(ctx, query, map[string]any{
			"id":             avatar.Id,
			"url":            avatar.Url,
			"isRelativePath": avatar.IsRelativePath,
		})
		if err != nil {
			return err
		}
		query = `
			INSERT INTO user_profiles(user_id, avatar_id)
			VALUES(:userId, :avatarId)
			ON CONFLICT (user_id)
			DO UPDATE SET
				avatar_id = :avatarId
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"userId":   userId,
			"avatarId": avatar.Id,
		})
		return err
	})
}
