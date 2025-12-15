package command

import (
	"context"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

var (
	ErrNoChange = errors.New("user no changes")
)

type UserUpdateInput struct {
	UserId uuid.UUID
	Name   string
	Email  string
}

func (c *UserCommand) UpdateUser(ctx context.Context, input UserUpdateInput) error {
	err := c.userRepository.Update(ctx, input.UserId, func(u *user.User) error {
		if u.Name == input.Name && u.Email.String() == input.Email {
			return ErrNoChange
		}
		err := u.SetName(input.Name)
		if err != nil {
			return err
		}
		err = u.SetEmail(input.Email)
		return err
	})
	if err == ErrNoChange {
		return nil
	}
	return nil
}
