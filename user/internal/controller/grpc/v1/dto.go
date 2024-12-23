package v1

import (
	"github.com/google/uuid"
	"github.com/meraiku/micro/user/internal/models"
	"github.com/meraiku/micro/user/pkg/auth_v1"
	"github.com/meraiku/micro/user/pkg/user_v1"
)

func FromEntity(user *models.User) *user_v1.User {
	return &user_v1.User{
		Id: user.ID.String(),
		Info: &user_v1.UserInfo{
			Name: user.Name,
		},
	}
}

func ToEntity(user *user_v1.User) (*models.User, error) {
	id, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, ErrInvalidID
	}
	return &models.User{
		ID:   id,
		Name: user.Info.Name,
	}, nil
}

func FromTokens(tokens *models.Tokens) *auth_v1.Tokens {
	return &auth_v1.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

func ToTokens(tokens *auth_v1.Tokens) *models.Tokens {
	return &models.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
