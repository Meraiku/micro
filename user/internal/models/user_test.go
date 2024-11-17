package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {

	tt := []struct {
		testName string
		name     string
		password string
		wantErr  error
	}{
		{
			testName: "empty name",
			name:     "",
			password: "password",
			wantErr:  ErrInvalidUserName,
		},
		{
			testName: "empty password",
			name:     "name",
			password: "",
			wantErr:  ErrEmptyPassword,
		},
		{
			testName: "valid user",
			name:     "name",
			password: "password",
			wantErr:  nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			_, err := NewUser(tc.name, tc.password)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestHashingPassword(t *testing.T) {
	pass := "password"

	user, _ := NewUser("name", pass)
	assert.NoError(t, user.HashPassword())

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(pass))
	assert.NoError(t, err)

	err = user.ValidatePassword(Password(pass))
	assert.NoError(t, err)

	err = user.ValidatePassword(Password("wrong"))
	assert.Error(t, err)
}
