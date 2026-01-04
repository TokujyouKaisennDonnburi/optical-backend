package query

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/github"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/api"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type MilestonesInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}

type MilestonesOutput struct {
	Title    string
	Progress int
	Open     int
	Close    int
}

func (g *GithubQuery) GetMilestone(ctx context.Context, input MilestonesInput) ([]MilestonesOutput, error) {
	// optionチェック
	options, err := g.optionRepository.FindsByCalendarId(ctx, input.CalendarId)
	if err != nil {
		return nil, err
	}
	hasOption := false
	for _, opt := range options {
		if opt.Id == option.OPTION_MILESTONE_STATUS {
			hasOption = true
			break
		}
	}
	if !hasOption {
		return nil, apperr.ForbiddenError("option not enabled")
	}
	milestones, err := g.githubRepository.GetMilestones(
		ctx,
		input.UserId,
		input.CalendarId,
		func(installationId string) (*github.Organization, error) {
			organization, err := g.stateRepository.GetOrganization(ctx, installationId)
			if err != nil {
				return nil, err
			}
			repositories, err := api.GetInstalledRepositories(ctx, organization.AccessToken)
			if err != nil {
				return nil, err
			}
			organization.SetRepositories(repositories)
			return organization, nil
		},
	)
	if err != nil {
		return nil, err
	}
	outputs := make([]MilestonesOutput, len(milestones))
	for i, milestone := range milestones {
		total := milestone.Open + milestone.Close
		progress := 0
		if total > 0 {
			progress = (milestone.Close * 100) / total
		}
		outputs[i] = MilestonesOutput{
			Title:    milestone.Title,
			Progress: progress,
			Open:     milestone.Open,
			Close:    milestone.Close,
		}
	}
	return outputs, nil
}
