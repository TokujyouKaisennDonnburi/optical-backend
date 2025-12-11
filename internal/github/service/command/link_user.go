package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
)

type GithubLinkUserInput struct {
	Code  string
	State string
}

func (c *GithubCommand) LinkUser(ctx context.Context, input GithubLinkUserInput) error {
	if input.Code == "" {
		return apperr.ValidationError("invalid code")
	}
	if input.State == "" {
		return apperr.ValidationError("invalid state")
	}
	return c.githubRepository.LinkUser(ctx, input.Code, input.State)
}
