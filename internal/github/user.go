package github

type User struct {
	Id        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string  `json:"id"`
	Url       string `json:"html_url"`
	AvatarUrl string `json:"avatar_url"`
}
