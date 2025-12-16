package gateway

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
)

func (r *OptionPsqlRepository) FindsByCalendarId(ctx context.Context, calendarId uuid.UUID) ([]option.Option, error) {
	query := `
		SELECT 
			options.id, options.name
		FROM calendar_options
		JOIN options
			ON calendar_options.option_id = options.id
		WHERE 
			calendar_options.calendar_id = $1
			AND options.deprecated = FALSE
	`
	optionModels := []psql.OptionModel{}
	err := r.db.SelectContext(ctx, &optionModels, query, calendarId)
	if err != nil {
		return nil, err
	}
	options := make([]option.Option, len(optionModels))
	for i, optionModel := range optionModels {
		options[i] = option.Option{
			Id:   optionModel.Id,
			Name: optionModel.Name,
		}
	}
	return options, nil
}

