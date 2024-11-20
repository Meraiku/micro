package containers

import (
	"github.com/meraiku/micro/user/internal/config"
	grpc_v1 "github.com/meraiku/micro/user/internal/controller/grpc/v1"
	rest_v1 "github.com/meraiku/micro/user/internal/controller/rest/v1"
	"github.com/meraiku/micro/user/internal/service/user"
)

type UserContainerGRPC struct {
	repo    *UserServiceRepos
	UserAPI *grpc_v1.GRPCUserService
}

func NewUserGRPC(r map[config.Repo]config.RepoType) (*UserContainerGRPC, error) {

	repos, err := NewUserServiceRepos(r)
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

type UserContainerREST struct {
	repo    *UserServiceRepos
	UserAPI *rest_v1.API
}

func NewUserREST(r map[config.Repo]config.RepoType) (*UserContainerREST, error) {

	repos, err := NewUserServiceRepos(r)
	if err != nil {
		return nil, err
	}

	userService := user.New(repos.user)

	api := rest_v1.New(userService)

	return &UserContainerREST{
		repo:    repos,
		UserAPI: api,
	}, nil
}
