package command

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/repository"
)

type GithubCommand struct {
	githubRepository repository.GithubRepository
}

func NewGithubCommand(githubRepository repository.GithubRepository) *GithubCommand {
	if githubRepository == nil {
		panic("githubRepository is nil")
	}
	return &GithubCommand{
		githubRepository: githubRepository,
	}
}
