package schedule

import (
	"errors"
	"unicode/utf8"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
)

const (
	MIN_SCHEDULE_NAME_LEN = 1
	MAX_SCHEDULE_NAME_LEN = 32
)

type Schedule struct {
	Id       uuid.UUID
	Name     string
	Calendar Calendar
	Members  []Member
	Options  []option.Option
}

func NewSchedule(name string, members []Member, options []option.Option) (*Schedule, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	nameLength := utf8.RuneCountInString(name)
	if nameLength < MIN_SCHEDULE_NAME_LEN || nameLength > MAX_SCHEDULE_NAME_LEN {
		return nil, errors.New("Schedule `name` is invalid")
	}
	return &Schedule{
		Id:       id,
		Name:     name,
		Calendar: Calendar{},
		Members:  members,
		Options:  options,
	}, nil
}

func (s *Schedule) AssignEvent(event Event) error {
	s.Calendar = s.Calendar.append(event)
	return nil
}
