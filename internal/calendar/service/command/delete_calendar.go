package command

import (
	"context"

	"github.com/google/uuid"
)

type CalendarDeleteInput struct {
	CalendarId uuid.UUID
	UserId     uuid.UUID
}

func (c *CalendarCommand) DeleteCalendar(ctx context.Context, input CalendarDeleteInput) error {
	return c.transactor.Transact(ctx, func(ctx context.Context) error {
		// 削除前にカレンダー情報を取得
		cal, err := c.calendarRepository.FindByUserCalendarId(ctx, input.UserId, input.CalendarId)
		if err != nil {
			return err
		}

		// メンバーへ通知
		_ = c.calendarNoticeService.NotifyCalendarDeleted(
			ctx,
			input.CalendarId,
			cal.Name,
			input.UserId,
		)

		// カレンダー削除
		return c.calendarRepository.Delete(ctx, input.CalendarId, input.UserId)
	})
}
