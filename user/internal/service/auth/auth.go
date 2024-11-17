package auth

import (
	"context"

	"github.com/meraiku/micro/user/internal/models"
)

type Repository interface {
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetTokens(ctx context.Context, user *models.User) (*models.Tokens, error) {
	return nil, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*models.Tokens, error) {
	return nil, nil
}
