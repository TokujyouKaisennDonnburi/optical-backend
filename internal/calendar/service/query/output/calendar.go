package output

import (
	"github.com/google/uuid"
)

type CalendarQueryOutput struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Color string    `json:"color"`
}
