package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type GithubLinkInput struct {
	UserId         uuid.UUID
	InstallationId string
}

func (c *GithubCommand) LinkUser(ctx context.Context, input GithubLinkInput) error {
	if input.InstallationId == "" {
		return apperr.ValidationError("invalid installationId")
	}
	err := c.githubRepository.LinkUser(ctx, input.UserId, input.InstallationId)
	return err
}
