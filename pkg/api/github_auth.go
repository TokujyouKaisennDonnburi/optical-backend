package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/auth"
)

type GithubAccessTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
}

type GithubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func PostOauthAccessToken(code string) (string, error) {
	client := http.Client{}
	requestBody := GithubAccessTokenRequest{
		Code:         code,
		ClientId:     auth.GetClientId(),
		ClientSecret: auth.GetClientSecret(),
		RedirectURI:  auth.GetGithubOauthRedirectURI(),
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(body),
	)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var respBody GithubAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", err
	}
	return respBody.AccessToken, nil
}
