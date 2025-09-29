package profile

import (
	"errors"
	"time"
)

type Profile struct {
	ID        string
	UserID    string
	PhotoURL  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrEmptyPhotoURL   = errors.New("写真URLが空です")
)

// NewProfile
func NewProfile(userID string) (*Profile, error) {
	now := time.Now()
	return &Profile{
		UserID:    userID,
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
