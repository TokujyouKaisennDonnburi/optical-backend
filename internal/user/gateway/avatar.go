package gateway

import (
	"context"
	"mime/multipart"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type AvatarPsqlAndMinioRepository struct {
	db          *sqlx.DB
	minioClient *minio.Client
	bucketName  string
}

func NewAvatarPsqlAndMinioRepository(db *sqlx.DB, minoiClient *minio.Client, bucketName string) *AvatarPsqlAndMinioRepository {
	if db == nil {
		panic("psql db is nil")
	}
	if minoiClient == nil {
		panic("minioClient is nil")
	}
	if bucketName == "" {
		panic("bucketName is empty")
	}
	return &AvatarPsqlAndMinioRepository{
		db:          db,
		minioClient: minoiClient,
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
	return r.bucketName + "/" + filePath + header.Filename, nil
}

func (r *AvatarPsqlAndMinioRepository) Save(ctx context.Context, avatar *user.Avatar) error {
	query := `
		INSERT INTO avatars(id, url)
		VALUES(:id, :url)
	`
	_, err := r.db.NamedExecContext(ctx, query, map[string]any{
		"id":  avatar.Id,
		"url": avatar.Url,
	})
	return err
}
