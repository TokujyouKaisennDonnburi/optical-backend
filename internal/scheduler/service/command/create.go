package command

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/google/uuid"
)

type SchedulerCreateInput struct {
	Title     string
	Memo      string
	StartTime time.Time
	EndTime   time.Time
	LimitTime time.Time
	IsAllDay  bool
}

func (c *SchedulerCommand) CreateScheduler(ctx context.Context, input SchedulerCreateInput)(*scheduler.Scheduler, error){
	// 予定時間を作成
	scheduledTime, err := calendar.NewScheduledTime(input.IsAllDay, input.StartTime, input.EndTime)
	if err != nil {
		return nil, err
	}
}

