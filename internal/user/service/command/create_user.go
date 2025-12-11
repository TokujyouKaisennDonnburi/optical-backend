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
	Id           uuid.UUID
	Name         string
	Email        string
	AccessToken  string
	RefreshToken string
}

func (c *UserCommand) CreateUser(ctx context.Context, input UserCreateInput) (*UserCreateOutput, error) {
	newUser, err := user.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}
	err = c.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	// アクセストークン発行
	accessToken, err := user.NewAccessToken(newUser.Id, newUser.Name)
	if err != nil {
		return nil, err
	}
	// リフレッシュトークン発行
	refreshToken, err := user.NewRefreshToken(newUser)
	if err != nil {
		return nil, err
	}
	// リフレッシュトークンを保存
	err = c.tokenRepository.AddToWhitelist(refreshToken)
	if err != nil {
		return nil, err
	}
	return &UserCreateOutput{
		Id:           newUser.Id,
		Name:         newUser.Name,
		Email:        newUser.Email,
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	}, nil
}
