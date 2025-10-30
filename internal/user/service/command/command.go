package command

import "github.com/TokujouKaisenDonburi/optical-backend/internal/user/repository"

type UserCommand struct {
	userRepository repository.UserRepository
}

func NewUserCommand(userRepository repository.UserRepository) *UserCommand {
	if userRepository == nil {
		panic("UserRepository is nil")
	}
	return &UserCommand{
		userRepository: userRepository,
	}
}
