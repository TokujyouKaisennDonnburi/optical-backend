package github

const (
	MILESTONES_STATE_OPEN  = "open"
	MILESTONES_STATE_CLOSE = "close"
)

type Milestones struct {
	Title string `json:"title"`
	Open  int    `json:"open"`
	Close int    `json:"close"`
}
