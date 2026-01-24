package command

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/storage"
	"github.com/google/uuid"
)

const (
	MAX_IMAGE_SIZE  = 20_000_000 // 20MB
	VALID_IMAGE_EXT = "png,jpg,jpeg"
)

type UploadAvatarInput struct {
	UserId uuid.UUID
	File   multipart.File
	Header *multipart.FileHeader
}

type UploadAvatarOutput struct {
	Id  uuid.UUID
	Url string
}

// 画像をストレージにアップロードしプロフィールに設定
func (c *UserCommand) UploadAvatar(ctx context.Context, input UploadAvatarInput) (*UploadAvatarOutput, error) {
	// ファイルサイズ確認
	if input.Header.Size > MAX_IMAGE_SIZE {
		return nil, apperr.ValidationError("invalid image size")
	}
	found := false
	// ファイルフォーマット確認
	for ext := range strings.SplitSeq(VALID_IMAGE_EXT, ",") {
		if strings.HasSuffix(input.Header.Filename, "."+ext) {
			found = true
			break
		}
	}
	if !found {
		return nil, apperr.ValidationError("invalid image ext")
	}
	// アバター情報を作成
	avatar, err := user.NewAvatar("", true)
	if err != nil {
		return nil, err
	}
	input.Header.Filename = avatar.Id.String() + ".png"
	// 画像をストレージにアップロード
	url, err := c.avatarRepository.Upload(ctx, input.File, input.Header)
	if err != nil {
		return nil, err
	}
	// URLを設定（ベースURLなしのパスを保存）
	avatar.SetUrl(url)
	// アバターをリポジトリに保存
	err = c.avatarRepository.Save(ctx, input.UserId, avatar)
	if err != nil {
		return nil, err
	}
	return &UploadAvatarOutput{
		Id:  avatar.Id,
		Url: storage.GetImageStorageBaseUrl() + "/" + avatar.Url,
	}, nil
}
