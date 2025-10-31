package command

import "github.com/TokujouKaisenDonburi/optical-backend/internal/user/repository"

type UserCommand struct {
	userRepository repository.UserRepository
	tokenRepository repository.TokenRepository
}

func NewUserCommand(userRepository repository.UserRepository, tokenRepository repository.TokenRepository) *UserCommand {
	if userRepository == nil {
		panic("UserRepository is nil")
	}
	if tokenRepository == nil {
		panic("TokenRepository is nil")
	}
	return &UserCommand{
		userRepository: userRepository,
		tokenRepository: tokenRepository,
	}
}
