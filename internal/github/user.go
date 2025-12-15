package github

type User struct {
	Id        int64  `json:"id"`
	Name      string `json:"login"`
	Email     string `json:"email"`
	Url       string `json:"html_url"`
	AvatarUrl string `json:"avatar_url"`
}
