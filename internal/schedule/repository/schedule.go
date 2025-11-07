package repository

  type ScheduleRepository interface {
      FindAll(ctx context.Context) ([]*schedule.Schedule, error)
  }

