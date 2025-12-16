package github

const (
	PULL_REQUEST_STATE_ALL   = "all"
	PULL_REQUEST_STATE_OPEN  = "open"
	PULL_REQUEST_STATE_CLOSE = "close"
)

type PullRequest struct {
	Id        int64  `json:"id"`
	Url       string `json:"html_url"`
	Title     string `json:"title"`
	State     string `json:"state"`
	Draft     bool   `json:"draft"`
	Number    int    `json:"number"`
	Assignees []User `json:"assignees"`
	Reviewers []User `json:"requested_reviewers"`
}
