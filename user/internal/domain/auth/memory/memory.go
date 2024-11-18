package memory

import (
	"context"
	"sync"

	"github.com/meraiku/micro/user/internal/domain/auth"
	"github.com/meraiku/micro/user/internal/models"
)

type Repository struct {
	store map[string]*models.Tokens
	mu    sync.RWMutex
}

func New() *Repository {
	return &Repository{
		store: make(map[string]*models.Tokens),
		mu:    sync.RWMutex{},
	}
}

func (r *Repository) StashTokens(ctx context.Context, userID string, tokens *models.Tokens) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[userID] = tokens

	return nil
}

func (r *Repository) GetTokens(ctx context.Context, userID string) (*models.Tokens, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tokens, ok := r.store[userID]
	if !ok {
		return nil, auth.ErrNoTokens
	}

	return tokens, nil
}
