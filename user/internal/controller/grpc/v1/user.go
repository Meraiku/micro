package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/meraiku/micro/user/internal/models"
	"github.com/meraiku/micro/user/pkg/user_v1"
)

type UserService interface {
	Get(ctx context.Context, id uuid.UUID) (*models.User, error)
	List(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GRPCServer struct {
	user_v1.UnimplementedUserV1Server
	userService UserService
}

func New(userService UserService) *GRPCServer {
	return &GRPCServer{
		userService: userService,
	}
}

func (s *GRPCServer) Get(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	return nil, nil
}

func (s *GRPCServer) List(ctx context.Context) (*user_v1.ListResponse, error) {
	return nil, nil
}

func (s *GRPCServer) Create(ctx context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	return nil, nil
}

func (s *GRPCServer) Update(ctx context.Context, req *user_v1.UpdateRequest) (*user_v1.UpdateResponse, error) {
	return nil, nil
}

func (s *GRPCServer) Delete(ctx context.Context, req *user_v1.DeleteRequest) error {
	return nil
}
