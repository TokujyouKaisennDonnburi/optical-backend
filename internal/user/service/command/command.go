package command

import "github.com/TokujouKaisenDonburi/optical-backend/internal/user/repository"

type UserCommand struct {
	userRepository   repository.UserRepository
	tokenRepository  repository.TokenRepository
	avatarRepository repository.AvatarRepository
}

func NewUserCommand(
	userRepository repository.UserRepository,
	tokenRepository repository.TokenRepository,
	avatarRepository repository.AvatarRepository,
) *UserCommand {
	if userRepository == nil {
		panic("userRepository is nil")
	}
	if tokenRepository == nil {
		panic("kokenRepository is nil")
	}
	if avatarRepository == nil {
		panic("avatarRepository is nil")
	}
	return &UserCommand{
		userRepository:   userRepository,
		tokenRepository:  tokenRepository,
		avatarRepository: avatarRepository,
	}
}
