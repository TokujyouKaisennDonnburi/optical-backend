package command

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
)

const (
	MAX_IMAGE_SIZE = 20_000_000 // 20MB
)

type SaveImageCommandInput struct {
	File   multipart.File
	Header *multipart.FileHeader
}

type SaveImageCommandOutput struct {
	Id uuid.UUID
}

// 画像をアップロードしてURLを保存する
func (c *CalendarCommand) SaveImage(ctx context.Context, input SaveImageCommandInput) (*SaveImageCommandOutput, error) {
	if input.Header.Size > MAX_IMAGE_SIZE {
		return nil, errors.New("Image size is invalid")
	}
	// 画像情報を作成
	image, err := calendar.NewImage("")
	if err != nil {
		return nil, err
	}
	input.Header.Filename = image.Id.String()+".png"
	// 画像をアップロードする
	url, err := c.imageRepository.Upload(ctx, input.File, input.Header)
	if err != nil {
		return nil, err
	}
	// URLを設定
	image.SetUrl(url)
	// 画像情報を保存
	err = c.imageRepository.Save(ctx, image)
	if err != nil {
		return nil, err
	}
	return &SaveImageCommandOutput{
		Id: image.Id,
	}, nil
}
