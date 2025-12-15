package user

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Image     Avatar
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewProfile
func NewProfile(userId uuid.UUID, imageUrl string) (*Profile, error) {
	now := time.Now()
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	img := NewAvatar(imageUrl)
	return &Profile{
		Id:        id,
		UserId:    userId,
		Image:     *img,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
