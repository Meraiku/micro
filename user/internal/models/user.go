package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInvalidUserName = errors.New("invalid user name")
)

type User struct {
	ID   uuid.UUID
	Name string
}

func NewUser(name string) (*User, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUserName, "empty user name")
	}
	return &User{
		ID:   uuid.New(),
		Name: name,
	}, nil
}
