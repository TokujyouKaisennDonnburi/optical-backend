package github

import "time"

type Organization struct {
	Id             int64        `json:"id"`
	Login          string       `json:"login"`
	Name           string       `json:"name"`
	InstallationId string       `json:"installation_id"`
	AccessToken    string       `json:"accessToken"`
	TokenExpiresAt time.Time    `json:"expires_at"`
	Repositories   []Repository `json:"repositories"`
}

func (org *Organization) SetRepositories(repositories []Repository) {
	org.Repositories = repositories
}
