package notice

// ドメインモデル
import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	MIN_NOTICE_TITLE_LENGTH   = 1
	MAX_NOTICE_TITLE_LENGTH   = 64
	MAX_NOTICE_CONTENT_LENGTH = 1024
)

type Notice struct {
	Id         uuid.UUID
	UserId     uuid.UUID
	EventId    uuid.NullUUID
	CalendarId uuid.NullUUID
	Title      string
	Content    string
	IsRead     bool
	CreatedAt  string
}

func NewNotice(userID uuid.UUID, eventID, calendarID uuid.NullUUID, title, content string) (*Notice, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	n := &Notice{
		Id:         id,
		UserId:     userID,
		EventId:    eventID,
		CalendarId: calendarID,
		IsRead:     false,
	}

	if err := n.SetTitle(title); err != nil {
		return nil, err
	}
	if err := n.SetContent(content); err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Notice) SetTitle(title string) error {
	titleLength := utf8.RuneCountInString(title)
	if titleLength < MIN_NOTICE_TITLE_LENGTH || titleLength > MAX_NOTICE_TITLE_LENGTH {
		return errors.New("Notice `title` length is invalid")
	}
	n.Title = title
	return nil
}

func (n *Notice) SetContent(content string) error {
	contentLength := utf8.RuneCountInString(content)
	if contentLength > MAX_NOTICE_CONTENT_LENGTH {
		return errors.New("Notice `content` length is invalid")
	}
	n.Content = content
	return nil
}
