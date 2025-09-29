package option

import (
	"errors"

	"github.com/google/uuid"
)

type Option struct {
	id   uuid.UUID
	name string
}

func NewOption(id uuid.UUID, name string) (*Option, error) {
	if id == uuid.Nil {
		return nil, errors.New("Option `id` is nil")
	}
	if name == "" {
		return nil, errors.New("Option `name` is nil")
	}
	return &Option{
		id:   id,
		name: name,
	}, nil
}
