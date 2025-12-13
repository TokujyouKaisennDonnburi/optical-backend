package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
)

type GithubUserResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	Url       string `json:"html_url"`
	AvatarUrl string `json:"avatar_url"`
}

func GetGithubUser(accessToken string) (*github.User, error) {
	client := http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("X-GitHub-Api-Version", GITHUB_API_VERSION)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get github user: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var respBody GithubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	email, err := GetGithubPrimaryEmail(accessToken)
	if err != nil {
		return nil, err
	}
	return &github.User{
		Id:        respBody.Id,
		Name:      respBody.Name,
		Email:     email,
		Url:       respBody.Url,
		AvatarUrl: respBody.AvatarUrl,
	}, nil
}
