package github

type User struct {
	Id        int64  `json:"id"`
	Login     string `json:"login"`
	Url       string `json:"html_url"`
	AvatarUrl string `json:"avatar_url"`
}
