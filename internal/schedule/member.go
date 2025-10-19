package schedule

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// スケジュールに参加するユーザー
type Member struct {
	UserId   uuid.UUID
	Name     string
	JoinedAt time.Time
}

func NewMember(userId uuid.UUID, name string, joinedAt time.Time) (*Member, error) {
	if userId == uuid.Nil {
		return nil, errors.New("Member `id` is nil")
	}
	if name == "" {
		return nil, errors.New("Member `name` is nil")
	}
	return &Member{
		UserId:   userId,
		Name:     name,
		JoinedAt: joinedAt,
	}, nil
}

