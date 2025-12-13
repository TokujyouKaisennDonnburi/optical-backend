package handler

import "github.com/TokujouKaisenDonburi/optical-backend/internal/github/service/command"

type GithubHandler struct {
	githubCommand *command.GithubCommand
}

func NewGithubHandler(
	githubCommand *command.GithubCommand,
) *GithubHandler {
	if githubCommand == nil {
		panic("githubCommand is nil")
	}
	return &GithubHandler{
		githubCommand: githubCommand,
	}
}
