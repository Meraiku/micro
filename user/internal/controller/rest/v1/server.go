package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/micro/user/internal/models"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidUserName = errors.New("invalid user name")
	ErrInvalidID       = errors.New("invalid id")
)

type UserService interface {
	Get(ctx context.Context, id uuid.UUID) (*models.User, error)
	List(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type API struct {
	userService UserService
}

func New(
	userService UserService,
) *API {
	return &API{
		userService: userService,
	}
}

func (api *API) Run(addr string) error {

	srv := &http.Server{
		Addr:         addr,
		Handler:      api.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return srv.ListenAndServe()
}
