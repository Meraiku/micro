package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUserName = errors.New("invalid user name")
	ErrEmptyPassword   = errors.New("empty password")
)

type User struct {
	ID       uuid.UUID
	Name     string
	Password Password
}

type Password []byte

func (pass Password) hash() (Password, error) {

	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	return Password(hash), nil
}

func NewUser(name, password string) (*User, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUserName, "empty user name")
	}

	if password == "" {
		return nil, ErrEmptyPassword
	}

	return &User{
		ID:       uuid.New(),
		Name:     name,
		Password: Password(password),
	}, nil
}

func (u *User) HashPassword() error {
	var err error

	u.Password, err = u.Password.hash()

	return err
}

func (u *User) ValidatePassword(password Password) error {
	return bcrypt.CompareHashAndPassword(u.Password, password)
}
