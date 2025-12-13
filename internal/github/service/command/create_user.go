package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
	"github.com/google/uuid"
)

type GithubCreateUserOutput struct {
	Url string
}

func (c *GithubCommand) CreateUser(ctx context.Context) (*GithubCreateUserOutput, error) {
	scopes := "read:user,user:email"
	state, err := GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	err = c.stateRepository.SaveOauthState(ctx, uuid.Max, state)
	if err != nil {
		return nil, err
	}
	url := "https://github.com/login/oauth/authorize?"
	url += "client_id=" + auth.GetClientId()
	url += "&redirect_uri=" + auth.GetGithubOauthRedirectURI()
	url += "&scope=" + scopes
	url += "&state=" + state
	return &GithubCreateUserOutput{
		Url: url,
	}, nil
}
