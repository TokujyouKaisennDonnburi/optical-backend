package command

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
)

type CreateGoogleStateInput struct {
	UserId uuid.UUID
}

type CreateGoogleStateOutput struct {
	Url string
}

func (c *UserCommand) CreateGoogleState(
	ctx context.Context,
	input CreateGoogleStateInput,
) (*CreateGoogleStateOutput, error) {
	stateCode, err := GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	err = c.oauthStateRepository.SaveOauthState(ctx, input.UserId, stateCode)
	if err != nil {
		return nil, err
	}
	clientId := c.oauthStateRepository.GetClientId()
	redirectUri := c.oauthStateRepository.GetRedirectUri()
	url := "https://accounts.google.com/o/oauth2/v2/auth?"
	url += "scope=email%20profile%20openid"
	url += "&include_granted_scopes=true"
	url += "&response_type=code"
	url += "&state=" + stateCode
	url += "&redirect_uri=" + redirectUri
	url += "&client_id=" + clientId
	return &CreateGoogleStateOutput{
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
