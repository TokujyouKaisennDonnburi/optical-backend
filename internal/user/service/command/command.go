package command

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user/repository"
)

type UserCommand struct {
	userRepository       repository.UserRepository
	tokenRepository      repository.TokenRepository
	avatarRepository     repository.AvatarRepository
	googleRepository     repository.GoogleRepository
	oauthStateRepository repository.OauthStateRepository
}

func NewUserCommand(
	userRepository repository.UserRepository,
	tokenRepository repository.TokenRepository,
	avatarRepository repository.AvatarRepository,
	googleRepository repository.GoogleRepository,
	oauthStateRepository repository.OauthStateRepository,
) *UserCommand {
	if userRepository == nil {
		panic("userRepository is nil")
	}
	if tokenRepository == nil {
		panic("tokenRepository is nil")
	}
	if avatarRepository == nil {
		panic("avatarRepository is nil")
	}
	if googleRepository == nil {
		panic("googleRepository is nil")
	}
	if oauthStateRepository == nil {
		panic("oauthStateRepository is nil")
	}
	return &UserCommand{
		userRepository:       userRepository,
		tokenRepository:      tokenRepository,
		avatarRepository:     avatarRepository,
		googleRepository:     googleRepository,
		oauthStateRepository: oauthStateRepository,
	}
}
