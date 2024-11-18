package v1

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/meraiku/micro/user/internal/models"
	"github.com/meraiku/micro/user/pkg/user_v1"
)

var (
	ErrInvalidID = errors.New("invalid user id")
)

var _ user_v1.UserV1Server = (*GRPCUserService)(nil)

type UserService interface {
	Get(ctx context.Context, id uuid.UUID) (*models.User, error)
	List(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type GRPCUserService struct {
	user_v1.UnimplementedUserV1Server
	userService UserService
}

func NewUserService(userService UserService) *GRPCUserService {
	return &GRPCUserService{
		userService: userService,
	}
}

func (s *GRPCUserService) Get(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, ErrInvalidID
	}

	user, err := s.userService.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	out := FromEntity(user)

	return &user_v1.GetResponse{User: out}, nil
}

func (s *GRPCUserService) List(ctx context.Context, _ *empty.Empty) (*user_v1.ListResponse, error) {

	users, err := s.userService.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*user_v1.User, len(users))
	for i, user := range users {
		out[i] = FromEntity(user)
	}

	return &user_v1.ListResponse{Users: out}, nil
}

func (s *GRPCUserService) Create(ctx context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {

	userInput := &models.User{
		Name: req.Info.Name,
	}

	user, err := s.userService.Create(ctx, userInput)
	if err != nil {
		return nil, err
	}

	out := FromEntity(user)

	return &user_v1.CreateResponse{User: out}, nil
}

func (s *GRPCUserService) Update(ctx context.Context, req *user_v1.UpdateRequest) (*user_v1.UpdateResponse, error) {

	userInput, err := ToEntity(req.User)
	if err != nil {
		return nil, err
	}

	user, err := s.userService.Update(ctx, userInput)
	if err != nil {
		return nil, err
	}

	out := FromEntity(user)

	return &user_v1.UpdateResponse{User: out}, nil
}

func (s *GRPCUserService) Delete(ctx context.Context, req *user_v1.DeleteRequest) (*empty.Empty, error) {

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, ErrInvalidID
	}

	if err := s.userService.Delete(ctx, id); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
