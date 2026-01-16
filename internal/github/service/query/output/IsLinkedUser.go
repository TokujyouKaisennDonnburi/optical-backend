package output

type IsLinkedUserQueryOutput struct {
	IsLinked    bool
	GithubId    string
	GithubName  string
	GithubEmail string
	IsSsoLogin  bool
	LinkedAt    string
}
