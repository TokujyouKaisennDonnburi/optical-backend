package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
)

func GetPullRequests(ctx context.Context, accessToken, owner, repos string) ([]github.PullRequest, error) {
	client := http.Client{}
	requestUrl := GITHUB_BASE_URL+"/repos/"+owner+"/"+repos+"/pulls"
	requestUrl += "?per_page=100"
	// インストール済みアプリ取得リクエスト
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		requestUrl,
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
		return nil, fmt.Errorf("failed to get github pullrequests: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	// WARNING: ドメインから直接参照
	var pullRequests []github.PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
		return nil, err
	}
	return pullRequests, nil
}
