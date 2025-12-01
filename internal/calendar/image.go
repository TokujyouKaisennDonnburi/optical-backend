package calendar

import (
	"net/url"

	"github.com/google/uuid"
)

type Image struct {
	Id    uuid.UUID
	Url   string
	Valid bool
}

func NewImage(imageUrl string) (*Image, error) {
	image := &Image{}
	// id生成
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	// id設定
	image.Id = id
	// URL設定
	image.SetUrl(imageUrl)
	return image, nil
}

// URLを設定する
func (img *Image) SetUrl(imageUrl string) {
	_, err := url.Parse(imageUrl)
	if imageUrl == "" || err != nil {
		img.Url = ""
		img.Valid = false
	} else {
		img.Url = imageUrl
		img.Valid = true
	}
}
