package command

import (
	"context"
	"errors"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
)

type TokenRefreshInput struct {
	RefreshToken string
}

type TokenRefreshOutput struct {
	AccessToken string
	ExpiresIn   time.Time
}

// リフレッシュトークンを使用してアクセストークンを発行する
func (u UserCommand) Refresh(ctx context.Context, input TokenRefreshInput) (*TokenRefreshOutput, error) {
	// リフレッシュトークンの情報を抽出
	refreshToken, err := user.DecodeRefreshToken(input.RefreshToken)
	if err != nil {
		return nil, err
	}
	if refreshToken.IsExpired() {
		return nil, errors.New("RefreshToken is expired")
	}
	// 有効なリフレッシュトークンか確認
	err = u.tokenRepository.IsWhitelisted(ctx, refreshToken.UserId, refreshToken.Id)
	if err != nil {
		return nil, err
	}
	// アクセストークン発行
	accessToken, err := user.NewAccessToken(refreshToken.UserId, refreshToken.UserName)
	if err != nil {
		return nil, err
	}
	return &TokenRefreshOutput{
		AccessToken: accessToken.Token,
		ExpiresIn:   accessToken.ExpiresIn,
	}, nil
}
