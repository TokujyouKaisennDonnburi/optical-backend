package command

import (
	"context"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type CalendarCreateArgs struct {
	UserId        uuid.UUID
	UserName      string
	CalendarName  string
	CalendarColor string
	ImageId       uuid.UUID
	OptionIds     []uuid.UUID
}

type CalendarCreateOutput struct {
	Id   uuid.UUID
	Name string
}

// カレンダーを新規作成する
func (s *CalendarCommand) CreateCalendar(ctx context.Context, args CalendarCreateArgs) (*CalendarCreateOutput, error) {
	var newCalendar *calendar.Calendar
	// カレンダーをリポジトリで作成
	err := s.calendarRepository.Create(ctx, args.UserId, args.ImageId, args.OptionIds, func(user *user.User, image *calendar.Image, options []option.Option) (*calendar.Calendar, error) {
		// オプションIDが全て正しいか確認
		if len(options) != len(args.OptionIds) {
			return nil, errors.New("invalid option ids")
		}
		// カレンダー作成者をメンバーとして作成
		member, err := calendar.NewMember(user.Id, user.Name)
		if err != nil {
			return nil, err
		}
		// カレンダーを作成
		newCalendar, err = calendar.NewCalendar(args.CalendarName, args.CalendarColor, *image, []calendar.Member{*member}, options)
		if err != nil {
			return nil, err
		}
		return newCalendar, nil
	})
	if err != nil {
		return nil, err
	}
	return &CalendarCreateOutput{
		Id:   newCalendar.Id,
		Name: newCalendar.Name,
	}, nil
}
