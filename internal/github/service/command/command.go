package command

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/repository"
)

type GithubCommand struct {
	stateRepository  repository.StateRepository
	githubRepository repository.GithubRepository
}

func NewGithubCommand(
	stateRepository repository.StateRepository,
	githubRepository repository.GithubRepository,
) *GithubCommand {
	if stateRepository == nil {
		panic("stateRepository is nil")
	}
	if githubRepository == nil {
		panic("githubRepository is nil")
	}
	return &GithubCommand{
		stateRepository:  stateRepository,
		githubRepository: githubRepository,
	}
}
