package schedule

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// スケジュールに参加するユーザー
type Member struct {
	id       uuid.UUID
	name     string
	joinedAt time.Time
}

func NewMember(id uuid.UUID, name string, joinedAt time.Time) (*Member, error) {
	if id == uuid.Nil {
		return nil, errors.New("Member `id` is nil")
	}
	if name == "" {
		return nil, errors.New("Member `name` is nil")
	}
	return &Member{
		id:       id,
		name:     name,
		joinedAt: joinedAt,
	}, nil
}
