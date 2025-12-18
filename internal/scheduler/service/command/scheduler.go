package command

import (
	"time"

	"github.com/google/uuid"
)
type Response struct {
	Id uuid.UUID
}


func SchedulerCreate (calendarId uuid.UUID, title, memo string, startTime, endTime time.Time, isAllDay bool)(*Response, error){
	
}
