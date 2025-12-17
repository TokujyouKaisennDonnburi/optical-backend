package command

import (
	"context"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type UserLoginInput struct {
	Email    string
	Password string
}

type UserLoginOutput struct {
	Id           uuid.UUID
	Name         string
	Email        string
	AccessToken  string
	RefreshToken string
}

// ユーザーにログインする
func (c *UserCommand) LoginUser(ctx context.Context, input UserLoginInput) (*UserLoginOutput, error) {
	// ユーザー取得
	loginUser, err := c.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if loginUser.IsDeleted() {
		return nil, errors.New("The user is deleted")
	}
	if !loginUser.VerifyPassword(input.Password) {
		return nil, errors.New("password is incorrect")
	}
	// アクセストークン発行
	accessToken, err := user.NewAccessToken(loginUser.Id, loginUser.Name)
	if err != nil {
		return nil, err
	}
	// リフレッシュトークン発行
	refreshToken, err := user.NewRefreshToken(loginUser)
	if err != nil {
		return nil, err
	}
	// リフレッシュトークンを保存
	err = c.tokenRepository.AddToWhitelist(refreshToken)
	if err != nil {
		return nil, err
	}
	return &UserLoginOutput{
		Id:           loginUser.Id,
		Name:         loginUser.Name,
		Email:        loginUser.Email.String(),
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	}, nil
}
