package user

import (
	"errors"
	"time"
	"net/url"

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
	ErrEmptyPhotoURL   = errors.New("photoURL is empty")
	ErrInvalidPhotoURL = errors.New("photoURL validate error")
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
	if err := validatePhotoURL(photoURL); err != nil {
		return err
	}
	p.PhotoURL = photoURL
	p.UpdatedAt = time.Now()
	return nil
}

// validate photoURL
func validatePhotoURL(photoURL string) error {
	if photoURL == "" {
		return ErrEmptyPhotoURL
	}
	_, err := url.Parse(photoURL)
	if err != nil {
		return ErrInvalidPhotoURL
	}
	return nil
}
