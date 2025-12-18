package repository

import (
	"context"

	"github.com/google/uuid"
)

type Scheduler interface {
	Create(
		ctx context.Context,
		id, calendarId uuid.UUID,
		title, memo string,
		startime, endtime uuid.Time,
		isAllDay bool,
	)error
}
