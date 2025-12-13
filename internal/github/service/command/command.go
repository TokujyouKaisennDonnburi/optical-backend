package command

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/repository"
	userRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/user/repository"
)

type GithubCommand struct {
	tokenRepository  userRepo.TokenRepository
	stateRepository  repository.StateRepository
	githubRepository repository.GithubRepository
}

func NewGithubCommand(
	tokenRepository userRepo.TokenRepository,
	stateRepository repository.StateRepository,
	githubRepository repository.GithubRepository,
) *GithubCommand {
	if tokenRepository == nil {
		panic("tokenRepository is nil")
	}
	if stateRepository == nil {
		panic("stateRepository is nil")
	}
	if githubRepository == nil {
		panic("githubRepository is nil")
	}
	return &GithubCommand{
		tokenRepository:  tokenRepository,
		stateRepository:  stateRepository,
		githubRepository: githubRepository,
	}
}
