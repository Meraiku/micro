package v1

import (
	"github.com/google/uuid"
	"github.com/meraiku/micro/user/internal/models"
)

type UserResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

func (req *CreateUserRequest) ToUser() (*models.User, error) {
	user, err := models.NewUser(req.Name)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (req *UpdateUserRequest) ToUser() (*models.User, error) {
	user, err := models.NewUser(req.Name)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func ToUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}
}