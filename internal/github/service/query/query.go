package query

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/repository"
	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
)

type GithubQuery struct {
	stateRepository  repository.StateRepository
	optionRepository optionRepo.OptionRepository
	githubRepository repository.GithubRepository
}

func NewGithubQuery(
	stateRepository repository.StateRepository,
	optionRepository optionRepo.OptionRepository,
	githubRepository repository.GithubRepository,
) *GithubQuery {
	if stateRepository == nil {
		panic("stateRepository is nil")
	}
	if optionRepository == nil {
		panic("optionRepository is nil")
	}
	if githubRepository == nil {
		panic("githubRepository is nil")
	}
	return &GithubQuery{
		stateRepository:  stateRepository,
		optionRepository: optionRepository,
		githubRepository: githubRepository,
	}
}
