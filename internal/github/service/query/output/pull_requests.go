package output

import "github.com/TokujouKaisenDonburi/optical-backend/internal/github"

type GithubPullRequestListQueryOutput struct {
	GithubId     int64
	Repository   github.Repository
	PullRequests []github.PullRequest
}
