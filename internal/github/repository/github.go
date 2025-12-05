package repository

import (
	"context"

	"github.com/google/uuid"
)

type GithubRepository interface {
	LinkUser(
		ctx context.Context,
		userId uuid.UUID,
		installationId string,
	) error
}
