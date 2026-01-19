package output

import "time"

type IsLinkedUserQueryOutput struct {
	IsLinked    bool
	GithubId    string
	GithubName  string
	GithubEmail string
	IsSsoLogin  bool
	LinkedAt    time.Time
}
