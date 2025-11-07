package command

import (
	"context"

	"github.com/google/uuid"
)

type ScheduleAllList struct{
	Id uuid.UUID
	Name string
}

func (c *ScheduleCommand) GetAllSchedules(ctx context.Context)([]ScheduleAllList, error) {
  schedules, err := c.scheduleRepository.FindAll(ctx)
  if err != nil{
	  return nil,err
  }
  result := make([]ScheduleAllList,len(schedules))
  for i,s := range schedules {
	  result[i] = ScheduleAllList{
		  Id: s.Id,
		  Name: s.Name,
	  }
  }
  return result,nil
}

