package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type UserCreateInput struct {
	Name     string
	Email    string
	Password string
}

type UserCreateOutput struct {
	Id uuid.UUID
}

func (c *UserCommand) CreateUser(ctx context.Context, input UserCreateInput) (*UserCreateOutput, error) {
	user, err := user.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}
	err = c.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &UserCreateOutput{
		Id: user.Id,
	}, nil
}
