package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type CalendarCreateInput struct {
	UserId        uuid.UUID
	UserName      string
	CalendarName  string
	CalendarColor string
	MemberEmails  []string
	ImageId       uuid.UUID
	OptionIds     []int32
}

type CalendarCreateOutput struct {
	Id   uuid.UUID
	Name string
}

// カレンダーを新規作成する
func (s *CalendarCommand) CreateCalendar(ctx context.Context, input CalendarCreateInput) (*CalendarCreateOutput, error) {
	var newCalendar *calendar.Calendar
	// カレンダーをリポジトリで作成
	err := s.calendarRepository.Create(ctx,
		input.ImageId,
		input.MemberEmails,
		input.OptionIds,
		func(image *calendar.Image, members []calendar.Member, options []option.Option) (*calendar.Calendar, error) {
			// オプションIDが全て正しいか確認
			if len(options) != len(input.OptionIds) {
				return nil, apperr.ValidationError("invalid option ids")
			}
			// 色作成
			color, err := calendar.NewColor(input.CalendarColor)
			if err != nil {
				return nil, err
			}
			// カレンダー作成者をメンバーとして作成
			member, err := calendar.NewMember(input.UserId, input.UserName)
			if err != nil {
				return nil, err
			}
			// 作成者のみ参加済みにする
			member.SetAsJoined()
			members = append(members, *member)
			// カレンダーを作成
			newCalendar, err = calendar.NewCalendar(input.CalendarName, color, *image, members, options)
			if err != nil {
				return nil, err
			}
			return newCalendar, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &CalendarCreateOutput{
		Id:   newCalendar.Id,
		Name: newCalendar.Name,
	}, nil
}
