package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/security"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CreateGoogleUserInput struct {
	Code  string
	State string
}

type CreateGoogleUserOutput struct {
	UserId       uuid.UUID
	AccessToken  string
	RefreshToken string
}

// Googleアカウントからユーザーを新規作成
func (c *UserCommand) CreateGoogleUser(
	ctx context.Context,
	input CreateGoogleUserInput,
) (*CreateGoogleUserOutput, error) {
	if input.Code == "" {
		return nil, apperr.ForbiddenError("invalid code")
	}
	if input.State == "" {
		return nil, apperr.ForbiddenError("invalid state")
	}
	// Stateチェック
	_, err := c.oauthStateRepository.GetOauthState(ctx, input.State)
	if err != nil {
		return nil, err
	}
	// Token取得
	token, err := c.googleRepository.GetTokenByCode(ctx, input.Code)
	if err != nil {
		return nil, err
	}
	// ユーザー情報取得
	googleUser, err := c.googleRepository.GetUserByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	appUser, err := c.googleRepository.FindUserByGoogleId(ctx, googleUser.Id)
	// ユーザーがない場合新規作成
	if err != nil {
		if !apperr.IsNotFound(err) {
			return nil, err
		}
		logrus.WithError(err).Info("creating new google user")
		password, err := security.GenerateRandomString(32)
		if err != nil {
			return nil, err
		}
		// ユーザー新規作成
		appUser, err = user.NewUser(googleUser.Name, googleUser.Email, password)
		if err != nil {
			return nil, err
		}
		avatar, err := user.NewAvatar(googleUser.AvatarUrl)
		if err != nil {
			return nil, err
		}
		// リポジトリ保存
		err = c.googleRepository.CreateUser(ctx, appUser, avatar, googleUser)
		if err != nil {
			return nil, err
		}
	}
	// アクセストークン発行
	accessToken, err := user.NewAccessToken(appUser.Id, appUser.Name)
	if err != nil {
		return nil, err
	}
	// リフレッシュトークン発行
	refreshToken, err := user.NewRefreshToken(appUser)
	if err != nil {
		return nil, err
	}
	// リフレッシュトークンを保存
	err = c.tokenRepository.AddToWhitelist(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return &CreateGoogleUserOutput{
		UserId:       appUser.Id,
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	}, nil
}
