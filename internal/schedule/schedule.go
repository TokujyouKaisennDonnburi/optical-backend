package schedule

import (
	"errors"
	"unicode/utf8"

	"github.com/TokujouKaisenDonburi/calendar-back/internal/option"
	"github.com/google/uuid"
)

const (
	MIN_SCHEDULE_NAME_LEN = 1
	MAX_SCHEDULE_NAME_LEN = 32
)

type Schedule struct {
	id       uuid.UUID
	name     string
	calender Calender
	members  []Member
	options  []option.Option
}

func NewSchedule(id uuid.UUID, name string, calender Calender, members []Member, options []option.Option) (*Schedule, error) {
	if id == uuid.Nil {
		return nil, errors.New("Schedule `id` is nil")
	}
	nameLength := utf8.RuneCountInString(name)
	if nameLength < MIN_SCHEDULE_NAME_LEN || nameLength > MAX_SCHEDULE_NAME_LEN {
		return nil, errors.New("Schedule `name` is invalid")
	}
	return &Schedule{
		id:       id,
		name:     name,
		calender: calender,
		members: members,
		options:  options,
	}, nil
}

func (s *Schedule) AssignEvent(event Event) error {
	s.calender = s.calender.append(event)
}
