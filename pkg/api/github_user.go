package api

import (
	"encoding/json"
	"net/http"
)

type GithubUserResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	Url       string `json:"html_url"`
	AvatarUrl string `json:"avatar_url"`
}

func GetGithubUser(accessToken string) (*GithubUserResponse, error) {
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
	defer resp.Body.Close()
	var respBody GithubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	return &respBody, nil
}

