package notice

// ドメインモデル
import (
	"github.com/google/uuid"
)

type Notice struct {
	Id         uuid.UUID
	UserId     uuid.UUID
	EventId    uuid.NullUUID
	CalendarId uuid.NullUUID
	Title      string
	Content    string
	IsRead     bool
}

func NewNotice(userID uuid.UUID, eventID, calendarID uuid.NullUUID, title, content string) (*Notice, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &Notice{
		Id:         id,
		UserId:     userID,
		EventId:    eventID,
		CalendarId: calendarID,
		Title:      title,
		Content:    content,
		IsRead:     false,
	}, nil
}
