package containers

import (
	grpc_v1 "github.com/meraiku/micro/user/internal/controller/grpc/v1"
	"github.com/meraiku/micro/user/internal/service/auth"
)

type AuthContainerGRPC struct {
	repo    *AuthServiceRepos
	AuthAPI *grpc_v1.GRPCAuthService
}

func NewAuthGRPC() (*AuthContainerGRPC, error) {

	repos, err := NewAuthServiceRepos()
	if err != nil {
		return nil, err
	}

	authService := auth.New(repos.user, repos.token)

	api := grpc_v1.NewAuthService(authService)

	return &AuthContainerGRPC{
		repo:    repos,
		AuthAPI: api,
	}, nil
}
