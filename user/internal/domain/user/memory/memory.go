package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	user_repo "github.com/meraiku/micro/user/internal/domain/user"
	"github.com/meraiku/micro/user/internal/models"
)

type Repository struct {
	store map[uuid.UUID]*models.User
	mu    sync.RWMutex
}

func New() *Repository {
	return &Repository{
		store: make(map[uuid.UUID]*models.User),
		mu:    sync.RWMutex{},
	}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.store[id]
	if !ok {
		return nil, user_repo.ErrUserNotFound
	}

	return user, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.store {
		if user.Name == username {
			return user, nil
		}
	}

	return nil, user_repo.ErrUserNotFound
}

func (r *Repository) List(ctx context.Context) ([]*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.store))
	for _, user := range r.store {
		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[user.ID]; ok {
		return nil, user_repo.ErrUserExists
	}

	for _, u := range r.store {
		if u.Name == user.Name {
			return nil, user_repo.ErrUserExists
		}
	}

	r.store[user.ID] = user

	return user, nil
}

func (r *Repository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[user.ID]; !ok {
		return nil, user_repo.ErrUserNotFound
	}

	r.store[user.ID] = user

	return user, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[id]; !ok {
		return user_repo.ErrUserNotFound
	}

	delete(r.store, id)

	return nil
}
