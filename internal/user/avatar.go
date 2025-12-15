package user

import (
	"net/url"

	"github.com/google/uuid"
)

type Avatar struct {
	Id    uuid.UUID
	Url   string
	Valid bool
}

func NewAvatar(url string) (*Avatar, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	img := &Avatar{}
	img.Id = id
	img.SetUrl(url)
	return img, nil
}

// URLを設定する
func (img *Avatar) SetUrl(imageUrl string) {
	_, err := url.Parse(imageUrl)
	if imageUrl == "" || err != nil {
		img.Url = ""
		img.Valid = false
	} else {
		img.Url = imageUrl
		img.Valid = true
	}
}
