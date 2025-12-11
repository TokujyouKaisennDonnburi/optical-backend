package command

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/google/uuid"
)

type GithubOauthStateInput struct {
	UserId uuid.UUID
}

type GithubOauthStateOutput struct {
	Url string
}

func (c *GithubCommand) CreateOauthState(
	ctx context.Context,
	input GithubOauthStateInput,
) (*GithubOauthStateOutput, error) {
	scopes := "read:user,real:email"
	state, err := GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	err = c.githubRepository.SaveOauthState(ctx, input.UserId, state)
	if err != nil {
		return nil, err
	}
	url := "https://github.com/login/oauth/authorize?"
	url += "client_id=" + auth.GetClientId()
	url += "&redirect_uri=" + auth.GetGithubOauthRedirectURI()
	url += "&scope=" + scopes
	url += "&state=" + state
	return &GithubOauthStateOutput{
		Url: url,
	}, nil
}

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}
