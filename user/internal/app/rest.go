package app

import (
	"os"

	v1 "github.com/meraiku/micro/user/internal/controller/rest/v1"
	"github.com/meraiku/micro/user/internal/domain/user/memory"
	"github.com/meraiku/micro/user/internal/service/user"
)

type restService struct {
	userRepo    user.Repository
	userService v1.UserService
}

func newRestService() *restService {
	return &restService{}
}

func (s *restService) Repo() user.Repository {
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

func (s *restService) Service() v1.UserService {
	if s.userService == nil {
		s.userService = user.New(s.Repo())
	}

	return s.userService
}
