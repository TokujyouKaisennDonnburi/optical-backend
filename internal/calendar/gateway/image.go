package gateway

import (
	"context"
	"mime/multipart"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type ImagePsqlAndMinioRepository struct {
	db          *sqlx.DB
	minioClient *minio.Client
	bucketName  string
}

func NewImagePsqlAndMinioRepository(db *sqlx.DB, minoiClient *minio.Client, bucketName string) *ImagePsqlAndMinioRepository {
	if db == nil {
		panic("psql db is nil")
	}
	if minoiClient == nil {
		panic("minioClient is nil")
	}
	if bucketName == "" {
		panic("bucketName is empty")
	}
	return &ImagePsqlAndMinioRepository{
		db:          db,
		minioClient: minoiClient,
		bucketName:  bucketName,
	}
}

func (r *ImagePsqlAndMinioRepository) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	filePath := "/calendars/images/"
	_, err := r.minioClient.PutObject(ctx, r.bucketName, filePath+header.Filename, file, header.Size, minio.PutObjectOptions{
		ContentType: "image/png",
	})
	if err != nil {
		return "", err
	}
	return header.Filename, nil
}

func (r *ImagePsqlAndMinioRepository) Save(ctx context.Context, image *calendar.Image) error {
	query := `
		INSERT INTO calendar_images(id, url)
		VALUES(:id, :url)
	`
	_, err := r.db.NamedExecContext(ctx, query, map[string]any{
		"id":  image.Id,
		"url": image.Url,
	})
	return err
}
