package command

import (
	"context"

	"github.com/google/uuid"
)

type CalendarUpdateInput struct {
	CalendarId    uuid.UUID
	UserId        uuid.UUID
	UserName      string
	CalendarName  string
	CalendarColor string
	MemberEmails  []string
	ImageId       uuid.UUID
	OptionIds     []int32
}

// カレンダーを更新する
func (c *CalendarCommand) UpdateCalendar(ctx context.Context, input CalendarCreateInput) error {

}
