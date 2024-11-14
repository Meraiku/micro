package app

import (
	"log"

	"github.com/meraiku/micro/user/internal/config"
	v1 "github.com/meraiku/micro/user/internal/controller/grpc/v1"
	"github.com/meraiku/micro/user/internal/service/user"
)

type grpcService struct {
	userRepo    user.Repository
	userService v1.UserService
	cfg         *config.GRPC
}

func newGRPCService() *grpcService {
	return &grpcService{}
}

func (s *grpcService) Config() *config.GRPC {
	if s.cfg == nil {
		cfg, err := config.NewGRPC()
		if err != nil {
			log.Fatalf("failed to load grpc config: %v", err)
		}

		s.cfg = cfg
	}

	return s.cfg
}

func (s *grpcService) Repo() user.Repository {
	if s.userRepo == nil {
		s.userRepo = setupUserRepository(s.Config().UserRepoType())
	}

	return s.userRepo
}

func (s *grpcService) Service() v1.UserService {
	if s.userService == nil {
		s.userService = user.New(s.Repo())
	}

	return s.userService
}
