package output

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar"
	"github.com/google/uuid"
)

type CalendarQueryOutput struct {
	Id    uuid.UUID
	Name  string
	Color string
	Image calendar.Image
}
