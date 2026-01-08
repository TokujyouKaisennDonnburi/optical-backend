package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/scheduler"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/google/uuid"
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

type schedulerCreateParams struct {
	Id         uuid.UUID `db:"id"`
	CalendarId uuid.UUID `db:"calendarId"`
	UserId     uuid.UUID `db:"userId"`
	Title      string    `db:"title"`
	Memo       string    `db:"memo"`
	LimitTime  time.Time `db:"limitTime"`
	IsAllDay   bool      `db:"isAllDay"`
}

type possibleDateParams struct {
	SchedulerId uuid.UUID `db:"schedulerId"`
	Date        time.Time `db:"date"`
	StartTime   time.Time `db:"startTime"`
	EndTime     time.Time `db:"endTime"`
}

// create
func (r *SchedulerPsqlRepository) CreateScheduler(
	ctx context.Context,
	id, calendarId, userId uuid.UUID,
	title, memo string,
	possibleDates []scheduler.PossibleDate,
	limitTime time.Time,
	isAllDay bool,
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		// scheduler
		sql := `
		INSERT INTO scheduler(id, calendar_id, user_id, title, memo, limit_time, is_allday)
		VALUES(:id, :calendarId, :userId, :title, :memo, :limitTime, :isAllDay)
	`
		schedulerParams := schedulerCreateParams{
			Id:         id,
			CalendarId: calendarId,
			UserId:     userId,
			Title:      title,
			Memo:       memo,
			LimitTime:  limitTime,
			IsAllDay:   isAllDay,
		}
		_, err := tx.NamedExecContext(ctx, sql, schedulerParams)
		if err != nil {
			return err
		}

		// possibleDate
		sql = `
		INSERT INTO scheduler_possible_date(scheduler_id, date, start_time, end_time)
		VALUES(:schedulerId, :date, :startTime, :endTime)
	`
		possibleDateParamsList := make([]possibleDateParams, len(possibleDates))
		for i, possibleDate := range possibleDates {
			possibleDateParamsList[i] = possibleDateParams{
				SchedulerId: id,
				Date:        possibleDate.Date,
				StartTime:   possibleDate.StartTime,
				EndTime:     possibleDate.EndTime,
			}
		}

		if len(possibleDateParamsList) == 0 {
			return nil
		}
		_, err = tx.NamedExecContext(ctx, sql, possibleDateParamsList)
		if err != nil {
			return err
		}
		return nil
	})
}
