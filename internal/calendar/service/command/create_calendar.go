package command

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

type CalendarCreateInput struct {
	UserId        uuid.UUID
	UserName      string
	CalendarName  string
	CalendarColor string
	ImageId       uuid.UUID
	OptionIds     []int32
}

type CalendarCreateOutput struct {
	Id   uuid.UUID
	Name string
}

// カレンダーを新規作成する
func (c *CalendarCommand) CreateCalendar(ctx context.Context, input CalendarCreateInput) (*CalendarCreateOutput, error) {
	var newCalendar *calendar.Calendar
	err := c.transactor.Transact(ctx, func(ctx context.Context) error {
		// 色作成
		color, err := calendar.NewColor(input.CalendarColor)
		if err != nil {
			return err
		}
		// オプション取得
		options, err := c.optionRepository.FindsByIds(ctx, input.OptionIds)
		if err != nil {
			return err
		}
		// オプションIDが全て正しいか確認
		if len(options) != len(input.OptionIds) {
			return apperr.ValidationError("invalid option ids")
		}
		image := &calendar.Image{Valid: false}
		// 画像を取得
		if input.ImageId != uuid.Nil {
			image, err = c.imageRepository.FindById(ctx, input.ImageId)
			if err != nil {
				return err
			}
		}
		// カレンダー作成者をメンバーとして作成
		member, err := calendar.NewMember(input.UserId, input.UserName)
		if err != nil {
			return err
		}
		// 作成者のみ参加済みにする
		member.SetAsJoined()
		members := []calendar.Member{*member}
		// カレンダーを作成
		newCalendar, err = calendar.NewCalendar(input.CalendarName, color, *image, members, options)
		if err != nil {
			return err
		}
		// カレンダーをリポジトリで作成
		err = c.calendarRepository.Create(ctx, newCalendar)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &CalendarCreateOutput{
		Id:   newCalendar.Id,
		Name: newCalendar.Name,
	}, nil
}
