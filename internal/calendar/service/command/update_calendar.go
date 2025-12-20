package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type CalendarUpdateInput struct {
	UserId        uuid.UUID
	CalendarId    uuid.UUID
	UserName      string
	CalendarName  string
	CalendarColor string
	MemberEmails  []string
	ImageId       uuid.UUID
	OptionIds     []int32
}

func (c *CalendarCommand) UpdateCalendar(ctx context.Context, input CalendarUpdateInput) error {
	err := c.calendarRepository.Update(
		ctx,
		input.CalendarId,
		input.ImageId,
		input.MemberEmails,
		input.OptionIds,
		func(existingCalendar *calendar.Calendar, image *calendar.Image, members []calendar.Member, options []option.Option) (*calendar.Calendar, error) {
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
			existingCalendar.SetImage(*image)
			if err := existingCalendar.SetMembers(members); err != nil {
				return nil, err
			}
			existingCalendar.SetOptions(options)
			return existingCalendar, nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}
