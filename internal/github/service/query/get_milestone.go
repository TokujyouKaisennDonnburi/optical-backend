package query

import (
	"context"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type MilestonesInput struct {
	UserId     uuid.UUID
	CalendarId uuid.UUID
}
type MilestonesOutput struct {
	Title    string
	Progress int8
	Open     int8
	Close    int8
}
type MilestonesItem struct {
	Owner string
	repo  string
	PAT   string
}

func (g GithubQuery) GetMilestone(ctx context.Context, input MilestonesInput) (*MileStonesOutput, error) {
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
		return nil, errors.New("")
	}
	if !hasOption {
		return nil, apperr.ForbiddenError("option is enabled")
	}
	outputs, err := g.githubRepository.GetMilestones(ctx, input.UserId, input.CalendarId,func(func(installationId string) (*github.Organization, error){
	}

}
