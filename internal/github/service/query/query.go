package query

import "github.com/TokujouKaisenDonburi/optical-backend/internal/github/repository"

type GithubQuery struct {
	stateRepository  repository.StateRepository
	githubRepository repository.GithubRepository
}

func NewGithubQuery(
	stateRepository repository.StateRepository,
	githubRepository repository.GithubRepository,
) *GithubQuery {
	if stateRepository == nil {
		panic("stateRepository is nil")
	}
	if githubRepository == nil {
		panic("githubRepository is nil")
	}
	return &GithubQuery{
		stateRepository:  stateRepository,
		githubRepository: githubRepository,
	}
}
