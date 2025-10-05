package option

import (
	"errors"

	"github.com/google/uuid"
)

type Option struct {
	Id   uuid.UUID
	Name string
}

func NewOption(name string) (*Option, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errors.New("Option `Name` is nil")
	}
	return &Option{
		Id:   id,
		Name: name,
	}, nil
}
