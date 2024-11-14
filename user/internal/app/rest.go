package app

import (
	"log"

	"github.com/meraiku/micro/user/internal/config"
	v1 "github.com/meraiku/micro/user/internal/controller/rest/v1"
	"github.com/meraiku/micro/user/internal/service/user"
)

type restService struct {
	userRepo    user.Repository
	userService v1.UserService
	cfg         *config.REST
}

func newRestService() *restService {
	return &restService{}
}

func (s *restService) Config() *config.REST {
	if s.cfg == nil {
		cfg, err := config.NewREST()
		if err != nil {
			log.Fatalf("failed to load rest config: %v", err)
		}

		s.cfg = cfg
	}

	return s.cfg
}

func (s *restService) Repo() user.Repository {
	if s.userRepo == nil {
		s.userRepo = setupUserRepository(s.Config().UserRepoType())
	}

	return s.userRepo
}

func (s *restService) Service() v1.UserService {
	if s.userService == nil {
		s.userService = user.New(s.Repo())
	}

	return s.userService
}
