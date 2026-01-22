package user

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Avatar    Avatar
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewProfile
func NewProfile(userId uuid.UUID, imageUrl string, isRelativePath bool) (*Profile, error) {
	now := time.Now().UTC()
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	avatar, err := NewAvatar(imageUrl, isRelativePath)
	if err != nil {
		return nil, err
	}
	return &Profile{
		Id:        id,
		UserId:    userId,
		Avatar:    *avatar,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
