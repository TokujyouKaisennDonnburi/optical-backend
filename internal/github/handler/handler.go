package handler

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/command"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/query"
)

type GithubHandler struct {
	githubQuery   *query.GithubQuery
	githubCommand *command.GithubCommand
}

func NewGithubHandler(
	githubQuery *query.GithubQuery,
	githubCommand *command.GithubCommand,
) *GithubHandler {
	if githubQuery == nil {
		panic("githubQuery is nil")
	}
	if githubCommand == nil {
		panic("githubCommand is nil")
	}
	return &GithubHandler{
		githubQuery:   githubQuery,
		githubCommand: githubCommand,
	}
}
