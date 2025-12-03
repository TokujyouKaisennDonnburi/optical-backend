package query

import (
	"github.com/TokujouKaisenDonburi/optical-backend/internal/calendar/repository"
)

// Event
type EventQuery struct {
	eventRepository repository.EventRepository
}

// コンストラクタ DI
func NewEventQuery(eventRepo repository.EventRepository) *EventQuery {
	return &EventQuery{
		eventRepository: eventRepo,
	}
}
