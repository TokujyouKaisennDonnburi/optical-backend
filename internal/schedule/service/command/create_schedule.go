package command

import (
	"context"
	"errors"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/schedule"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/user"
	"github.com/google/uuid"
)

type ScheduleCreateArgs struct {
	UserId       uuid.UUID
	UserName     string
	ScheduleName string
	OptionIds    []uuid.UUID
}

type ScheduleCreateOutput struct {
	Id   uuid.UUID
	Name string
}

// スケジュールを新規作成する
func (s *ScheduleCommand) CreateSchedule(ctx context.Context, args ScheduleCreateArgs) (*ScheduleCreateOutput, error) {
	var newSchedule *schedule.Schedule
	// スケジュールをリポジトリで作成
	err := s.scheduleRepository.Create(ctx, args.UserId, args.OptionIds, func(user *user.User, options []option.Option) (*schedule.Schedule, error) {
		// オプションIDが全て正しいか確認
		if len(options) != len(args.OptionIds) {
			return nil, errors.New("invalid option ids")
		}
		// スケジュール作成者をメンバーとして作成
		member, err := schedule.NewMember(user.Id, user.Name)
		if err != nil {
			return nil, err
		}
		// スケジュールを作成
		newSchedule, err = schedule.NewSchedule(args.ScheduleName, []schedule.Member{*member}, options)
		if err != nil {
			return nil, err
		}
		return newSchedule, nil
	})
	if err != nil {
		return nil, err
	}
	return &ScheduleCreateOutput{
		Id:   newSchedule.Id,
		Name: newSchedule.Name,
	}, nil
}
