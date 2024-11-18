package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/meraiku/micro/user/internal/models"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	List(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo Repository
}

func New(
	repo Repository,
) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) List(ctx context.Context) ([]*models.User, error) {

	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) Create(ctx context.Context, user *models.User) (*models.User, error) {

	out, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) Update(ctx context.Context, user *models.User) (*models.User, error) {

	out, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
