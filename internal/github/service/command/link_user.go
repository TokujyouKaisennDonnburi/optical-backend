package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type GithubLinkUserInput struct {
	Code  string
	State string
}

type GithubLinkUserOutput struct {
	UserId       uuid.UUID
	AccessToken  string
	RefreshToken string
}

func (c *GithubCommand) LinkUser(ctx context.Context, input GithubLinkUserInput) (*GithubLinkUserOutput, error) {
	if input.Code == "" {
		return nil, apperr.ValidationError("invalid code")
	}
	if input.State == "" {
		return nil, apperr.ValidationError("invalid state")
	}
	userId, err := c.stateRepository.GetOauthState(ctx, input.State)
	if err != nil {
		return nil, apperr.ForbiddenError("invalid state")
	}
	// uuidがMaxの場合ユーザーを新規作成
	if userId == uuid.Max {
		newUser, err := c.githubRepository.CreateUser(ctx, input.Code)
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
		return &GithubLinkUserOutput{
			UserId:       newUser.Id,
			AccessToken:  accessToken.Token,
			RefreshToken: refreshToken.Token,
		}, nil
	} else {
		err = c.githubRepository.LinkUser(ctx, userId, input.Code)
		return &GithubLinkUserOutput{
			UserId: userId,
		}, err
	}
}
