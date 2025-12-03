package output

import (
	"github.com/google/uuid"
)

type CalendarQueryOutput struct {
	Id    uuid.UUID
	Name  string
	Color string
}
