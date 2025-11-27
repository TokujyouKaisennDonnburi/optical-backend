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
	if imageUrl == "" {
		return &Image{
			Valid: false,
		}, nil
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	_, err = url.Parse(imageUrl)
	if err != nil {
		return nil, err
	}
	return &Image{
		Id:    id,
		Url:   imageUrl,
		Valid: true,
	}, nil
}
