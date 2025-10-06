package user

import (
	"errors"
	"time"

	"github.com/google/uuid"

)

type Profile struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	PhotoURL  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrEmptyPhotoURL   = errors.New("写真URLが空です")
)

// NewProfile
func NewProfile(userId uuid.UUID) (*Profile, error) {
	now := time.Now()
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &Profile{
		Id:        id,
		UserId:    userId,
		PhotoURL:  "",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// update photo
func (p *Profile) UpdatePhotoURL(photoURL string) error {
	if photoURL == "" {
		return ErrEmptyPhotoURL
	}
	p.PhotoURL = photoURL
	p.UpdatedAt = time.Now()
	return nil
}

// validate photoURL
// TODO #46
