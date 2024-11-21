package memory

import (
	"context"
	"sync"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/internal/models"
	"github.com/meraiku/micro/user/internal/service/auth"
)

type Repository struct {
	store map[string]models.Tokens
	mu    sync.RWMutex
}

func New() *Repository {
	return &Repository{
		store: make(map[string]models.Tokens),
		mu:    sync.RWMutex{},
	}
}

func (r *Repository) StashTokens(ctx context.Context, userID string, tokens *models.Tokens) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[userID] = *tokens

	logging.L(ctx).Debug(
		"tokens stashed",
		logging.String("user_id", userID),
		logging.Any("store", r.store),
	)

	return nil
}

func (r *Repository) GetTokens(ctx context.Context, userID string) (*models.Tokens, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	logging.L(ctx).Debug(
		"getting tokens",
		logging.String("user_id", userID),
		logging.Any("store", r.store),
	)

	tokens, ok := r.store[userID]
	if !ok {
		return nil, auth.ErrNoTokens
	}

	return &tokens, nil
}
