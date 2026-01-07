package gateway

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type SchedulerPsqlRepository struct {
	db *sqlx.DB
}

func NewSchedulerPsqlRepository(db *sqlx.DB) *SchedulerPsqlRepository {
	if db == nil {
		panic("db is nil")
	}
	return &SchedulerPsqlRepository{
		db: db,
	}
}

// create
func (r *SchedulerPsqlRepository) Create(ctx context.Context, id, calendarId, userId uuid.UUID, title, memo string, startTime, endTime, limitTime time.Time, isAllDay bool) error{
	// create scheduler
	sql := `
		INSERT INTO scheduler(id,calendarId,userId,title,memo,start_time,end_time,limit_time,is_allday)
		VALUES(:id, :calendarId, :userId, :title, memo, startTime, endTime, limitTime, isAllDay)
	`
	_, err := 
}
