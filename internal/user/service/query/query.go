package query

import "github.com/TokujouKaisenDonburi/optical-backend/internal/user/repository"

type UserQuery struct {
	userRepository repository.UserRepository
}

func NewUserQuery(userRepository repository.UserRepository) *UserQuery {
	if userRepository == nil {
		panic("UserRepository is nil")
	}
	return &UserQuery{
		userRepository: userRepository,
	}
}
