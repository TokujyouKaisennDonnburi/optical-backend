package todo

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	MIN_LIST_NAME_LENGTH = 1
	MAX_LIST_NAME_LENGTH = 32
)

type List struct {
	Id         uuid.UUID
	UserId     uuid.UUID
	CalendarId uuid.UUID
	Name       string
	Items      []Item
}

func NewList(userId, calendarId uuid.UUID, name string) (*List, error) {
	if userId == uuid.Nil {
		return nil, errors.New("todo list `userId` is nil")
	}
	if calendarId == uuid.Nil {
		return nil, errors.New("todo list `calendarId` is nil")
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	todoList := &List{
		Id:         id,
		UserId:     userId,
		CalendarId: calendarId,
		Items:      []Item{},
	}
	err = todoList.SetName(name)
	if err != nil {
		return nil, err
	}
	return todoList, nil
}

func (l *List) SetName(name string) error {
	length := utf8.RuneCountInString(name)
	if length < MIN_LIST_NAME_LENGTH || length > MAX_LIST_NAME_LENGTH {
		return errors.New("todo list `name` is invalid")
	}
	l.Name = name
	return nil
}

func (l *List) SetItems(items []Item) {
	l.Items = items
}
