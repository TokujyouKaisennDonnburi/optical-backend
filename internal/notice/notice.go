package notice

// ドメインモデル
import (
	"time"

	"github.com/google/uuid"
)

type Notice struct {
	Id        uuid.UUID
	EventId   uuid.UUID
	Title     string
	Content   string
	IsRead    bool
	CreatedAt time.Time
}

func NewNotice(eventID uuid.UUID, title, content string) (*Notice, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &Notice{
		Id:        id,
		EventId:   eventID,
		Title:     title,
		Content:   content,
		IsRead:    false,
		CreatedAt: time.Now(),
	}, nil
}
