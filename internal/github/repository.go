package github

type Repository struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	FullName  string `json:"fullName"`
}
