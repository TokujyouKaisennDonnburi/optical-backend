package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
)

type GithubMilestoneResponse struct {
	Title        string `json:"title"`
	OpenIssues   int    `json:"open_issues"`
	ClosedIssues int    `json:"closed_issues"`
}

func GetMilestones(ctx context.Context, accessToken, owner, repo, state string) ([]github.Milestones, error) {
	client := http.Client{}
	requestUrl := GITHUB_BASE_URL + "/repos/" + owner + "/" + repo + "/milestones"
	requestUrl += "?per_page=100"
	requestUrl += "&state=" + state
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
		return nil, fmt.Errorf("failed to get github milestones: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var respBody []GithubMilestoneResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	milestones := make([]github.Milestones, len(respBody))
	for i, milestone := range respBody {
		milestones[i] = github.Milestones{
			Title: milestone.Title,
			Open:  milestone.OpenIssues,
			Close: milestone.ClosedIssues,
		}
	}
	return milestones, nil
}
