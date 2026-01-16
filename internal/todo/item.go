package todo

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	MIN_ITEM_NAME_LENGTH = 1
	MAX_ITEM_NAME_LENGTH = 32
)

type Item struct {
	Id     uuid.UUID
	ListId uuid.UUID
	UserId uuid.UUID
	Name   string
	IsDone bool
}

func NewItem(listId, userId uuid.UUID, name string) (*Item, error) {
	if userId == uuid.Nil {
		return nil, errors.New("todo item `userId` is nil")
	}
	if listId == uuid.Nil {
		return nil, errors.New("todo item `listId` is nil")
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	todoItem := &Item{
		Id:     id,
		ListId: listId,
		UserId: userId,
		IsDone: false,
	}
	err = todoItem.SetName(name)
	if err != nil {
		return nil, err
	}
	return todoItem, nil
}

func (l *Item) SetName(name string) error {
	length := utf8.RuneCountInString(name)
	if length < MIN_ITEM_NAME_LENGTH || length > MAX_ITEM_NAME_LENGTH {
		return errors.New("todo item `name` is invalid")
	}
	l.Name = name
	return nil
}

func (l *Item) SetDone(done bool) {
	l.IsDone = done
}
