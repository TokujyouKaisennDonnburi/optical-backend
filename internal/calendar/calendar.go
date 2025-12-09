package calendar

import (
	"errors"
	"unicode/utf8"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/option"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/google/uuid"
)

const (
	MIN_CALENDAR_NAME_LEN = 1
	MAX_CALENDAR_NAME_LEN = 32
)

type Calendar struct {
	Id      uuid.UUID
	Name    string
	Color   Color
	Image   Image
	Members []Member
	Options []option.Option
}

func NewCalendar(name string, color Color, image Image, members []Member, options []option.Option) (*Calendar, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	nameLength := utf8.RuneCountInString(name)
	if nameLength < MIN_CALENDAR_NAME_LEN || nameLength > MAX_CALENDAR_NAME_LEN {
		return nil, errors.New("Calendar `name` is invalid")
	}
	if len(members) == 0 {
		return nil, apperr.ValidationError("Calendar `members` is empty")
	}
	return &Calendar{
		Id:      id,
		Name:    name,
		Color:   color,
		Image:   image,
		Members: members,
		Options: options,
	}, nil
}
