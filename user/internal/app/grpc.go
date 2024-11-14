package app

import (
	"os"

	v1 "github.com/meraiku/micro/user/internal/controller/grpc/v1"
	"github.com/meraiku/micro/user/internal/domain/user/memory"
	"github.com/meraiku/micro/user/internal/service/user"
)

type grpcService struct {
	userRepo    user.Repository
	userService v1.UserService
}

func newGRPCService() *grpcService {
	return &grpcService{}
}

func (s *grpcService) Repo() user.Repository {
	if s.userRepo == nil {
		var repo user.Repository

		switch os.Getenv("USER_REPO") {
		default:
			memoryRepo := memory.New()
			repo = memoryRepo
		}

		s.userRepo = repo
	}

	return s.userRepo
}

func (s *grpcService) Service() v1.UserService {
	if s.userService == nil {
		s.userService = user.New(s.Repo())
	}

	return s.userService
}
