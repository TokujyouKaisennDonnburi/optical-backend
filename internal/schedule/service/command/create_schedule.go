package command

import (
	"context"
	"errors"
	"time"

	optionRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/option/repository"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/schedule"
	scheduleRepo "github.com/TokujouKaisenDonburi/optical-backend/internal/schedule/repository"
	"github.com/google/uuid"
)

type CreateSchedule struct {
	scheduleRepository scheduleRepo.ScheduleRepository
	optionRepository   optionRepo.OptionRepository
}

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
func (s *CreateSchedule) CreateSchedule(ctx context.Context, args ScheduleCreateArgs) (*ScheduleCreateOutput, error) {
	// オプション一覧取得
	options, err := s.optionRepository.FindByIds(ctx, args.OptionIds)
	if err != nil {
		return nil, err
	}
	// オプションIDが全て正しいか確認
	if len(options) != len(args.OptionIds) {
		return nil, errors.New("invalid option ids")
	}
	// スケジュール作成者をメンバーとして作成
	member, err := schedule.NewMember(args.UserId, args.UserName, time.Time{})
	if err != nil {
		return nil, err
	}
	// スケジュールを作成
	schedule, err := schedule.NewSchedule(args.ScheduleName, []schedule.Member{*member}, options)
	if err != nil {
		return nil, err
	}
	// スケジュールをリポジトリで作成
	err = s.scheduleRepository.Create(ctx, schedule)
	if err != nil {
		return nil, err
	}
	return &ScheduleCreateOutput{
		Id:   schedule.Id,
		Name: schedule.Name,
	}, nil
}
