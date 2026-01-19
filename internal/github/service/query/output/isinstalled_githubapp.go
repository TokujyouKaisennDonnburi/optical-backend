package output

import "time"

type IsInstalledGithubAppQueryOutput struct {
	IsInstalled    bool
	GithubId       string
	GithubName     string
	InstallationId string
	InstalledAt    time.Time
}
