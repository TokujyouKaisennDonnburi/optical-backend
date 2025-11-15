package calendar

import (
	"errors"
	"unicode/utf8"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/google/uuid"
)

const (
	MIN_CALENDAR_NAME_LEN = 1
	MAX_CALENDAR_NAME_LEN = 32
)

type Calendar struct {
	Id       uuid.UUID
	Name     string
	Schedules Schedules
	Members  []Member
	Options  []option.Option
}

func NewCalendar(name string, members []Member, options []option.Option) (*Calendar, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	nameLength := utf8.RuneCountInString(name)
	if nameLength < MIN_CALENDAR_NAME_LEN || nameLength > MAX_CALENDAR_NAME_LEN {
		return nil, errors.New("Schedule `name` is invalid")
	}
	return &Calendar{
		Id:       id,
		Name:     name,
		Schedules: Schedules{},
		Members:  members,
		Options:  options,
	}, nil
}

func (s *Calendar) AssignEvent(event Event) error {
	s.Schedules = s.Schedules.append(event)
	return nil
}
