package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type GithubEmailResponse struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

func GetGithubPrimaryEmail(accessToken string) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user/emails",
		nil,
	)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("X-GitHub-Api-Version", GITHUB_API_VERSION)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get github user emails: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var respBody []GithubEmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", err
	}
	for _, email := range respBody {
		if email.Primary {
			return email.Email, nil
		}
	}
	return "", errors.New("primary email not found")
}
