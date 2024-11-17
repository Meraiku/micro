package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/micro/user/internal/models"
	"github.com/meraiku/micro/user/pkg/tokens"
)

var (
	ErrIncorrectPassword = errors.New("incorrect password")
)

type Repository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Login(ctx context.Context, user *models.User) (*models.Tokens, error) {

	u, err := s.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := u.ValidatePassword(user.Password); err != nil {
		return nil, ErrIncorrectPassword
	}

	access, err := tokens.GenerateJWT(u.ID.String(), 60*time.Minute, []byte("secret"))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refresh, err := tokens.GenerateJWT(u.ID.String(), 24*time.Hour, []byte("secret"))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokens, err := models.NewTokens(access, refresh)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}

	return tokens, nil
}

func (s *Service) Register(ctx context.Context, user *models.User) (*models.User, error) {

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	u, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return u, nil
}

func (s *Service) GetTokens(ctx context.Context, user *models.User) (*models.Tokens, error) {
	return nil, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*models.Tokens, error) {
	return nil, nil
}
