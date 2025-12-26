package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type CalendarUpdateInput struct {
	UserId        uuid.UUID
	CalendarId    uuid.UUID
	CalendarName  string
	CalendarColor string
	OptionIds     []int32
}

func (c *CalendarCommand) UpdateCalendar(ctx context.Context, input CalendarUpdateInput) error {

	// オプション取得
	options, err := c.optionRepository.FindOptionsByIds(ctx, input.OptionIds)
	if err != nil {
		return err
	}

	err = c.calendarRepository.Update(
		ctx,
		input.UserId,
		input.CalendarId,
		func(existingCalendar *calendar.Calendar) (*calendar.Calendar, error) {
			if len(options) != len(input.OptionIds) {
				return nil, apperr.ValidationError("invalid option ids")
			}
			color, err := calendar.NewColor(input.CalendarColor)
			if err != nil {
				return nil, err
			}
			if err := existingCalendar.SetName(input.CalendarName); err != nil {
				return nil, err
			}
			existingCalendar.SetColor(color)
			existingCalendar.SetOptions(options)
			return existingCalendar, nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}
