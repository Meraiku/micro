package v1

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"

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
	us UserService,
) *API {
	return &API{
		userService: us,
	}
}

func (api *API) Run() error {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	host := os.Getenv("HOST")

	addr := net.JoinHostPort(host, port)

	srv := &http.Server{
		Addr:    addr,
		Handler: api.routes(),
	}

	return srv.ListenAndServe()
}
