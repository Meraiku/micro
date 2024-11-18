package containers

import (
	grpc_v1 "github.com/meraiku/micro/user/internal/controller/grpc/v1"
	"github.com/meraiku/micro/user/internal/service/user"
)

type UserContainerGRPC struct {
	repo    *UserServiceRepos
	UserAPI *grpc_v1.GRPCUserService
}

func NewUserGRPC() (*UserContainerGRPC, error) {

	repos, err := NewUserServiceRepos()
	if err != nil {
		return nil, err
	}

	userService := user.New(repos.user)

	api := grpc_v1.NewUserService(userService)

	return &UserContainerGRPC{
		repo:    repos,
		UserAPI: api,
	}, nil
}
