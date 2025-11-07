package command

import(
"github.com/TokujouKaisenDonburi/optical-backend/internal/schedule/repository"
)

type ScheduleCommand struct {
  scheduleRepository repository.ScheduleRepository
}

func NewScheduleCommand(scheduleRepository repository.ScheduleRepository) *ScheduleCommand{
  if scheduleRepository == nil{
	  panic("ScheduleRepository is nil")
  }
  return &ScheduleCommand{
	  scheduleRepository: scheduleRepository,
  }
}
