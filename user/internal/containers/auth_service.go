package containers

import (
	"github.com/meraiku/micro/user/internal/config"
	grpc_v1 "github.com/meraiku/micro/user/internal/controller/grpc/v1"
	"github.com/meraiku/micro/user/internal/service/auth"
)

type AuthContainerGRPC struct {
	repo    *AuthServiceRepos
	AuthAPI *grpc_v1.GRPCAuthService
}

func NewAuthGRPC(cfg *config.Config) (*AuthContainerGRPC, error) {

	repos, err := NewAuthServiceRepos(cfg.Repos)
	if err != nil {
		return nil, err
	}

	authService, err := auth.New(cfg, repos.user, repos.token)
	if err != nil {
		return nil, err
	}

	api := grpc_v1.NewAuthService(authService)

	return &AuthContainerGRPC{
		repo:    repos,
		AuthAPI: api,
	}, nil
}
