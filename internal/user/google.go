package user

import "errors"

type GoogleUser struct {
	Id        string
	Name      string
	Email     string
	AvatarUrl string
}

func NewGoogleUser(id, name, email, avatarUrl string) (*GoogleUser, error) {
	if id == "" {
		return nil, errors.New("invalid id")
	}
	if name == "" {
		return nil, errors.New("invalid name")
	}
	if email == "" {
		return nil, errors.New("invalid email")
	}
	return &GoogleUser{
		Id:        id,
		Name:      name,
		Email:     email,
		AvatarUrl: avatarUrl,
	}, nil
}
