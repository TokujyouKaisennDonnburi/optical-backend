package profile

import (
	"errors"
	"time"

	"github.com/google/uuid"

)

type Profile struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PhotoURL  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrEmptyPhotoURL   = errors.New("写真URLが空です")
)

// NewProfile
func NewProfile(userID uuid.UUID) *Profile {
	now := time.Now()
	return &Profile{
		ID:        uuid.New(),
		UserID:    userID,
		PhotoURL:  "",
		CreatedAt: now,
		UpdatedAt: now,
	}
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
