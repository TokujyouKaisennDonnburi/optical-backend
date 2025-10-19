package schedule

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// スケジュールに参加するユーザー
type Member struct {
	Id       uuid.UUID
	Name     string
	JoinedAt time.Time
}

func NewMember(id uuid.UUID, name string, joinedAt time.Time) (*Member, error) {
	if id == uuid.Nil {
		return nil, errors.New("Member `id` is nil")
	}
	if name == "" {
		return nil, errors.New("Member `name` is nil")
	}
	return &Member{
		Id:       id,
		Name:     name,
		JoinedAt: joinedAt,
	}, nil
}
