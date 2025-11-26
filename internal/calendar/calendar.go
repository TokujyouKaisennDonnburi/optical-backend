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
	Id        uuid.UUID
	Name      string
	Color     string
	Members   []Member
	Options   []option.Option
}

func NewCalendar(name, color string, members []Member, options []option.Option) (*Calendar, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	nameLength := utf8.RuneCountInString(name)
	if nameLength < MIN_CALENDAR_NAME_LEN || nameLength > MAX_CALENDAR_NAME_LEN {
		return nil, errors.New("Calendar `name` is invalid")
	}
	if utf8.RuneCountInString(color) != 6 {
		return nil, errors.New("Calendar `color` is invalid")
	}
	return &Calendar{
		Id:        id,
		Name:      name,
		Color:     color,
		Members:   members,
		Options:   options,
	}, nil
}
