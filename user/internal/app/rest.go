package app

import (
	"context"
	"log"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/internal/config"
	v1 "github.com/meraiku/micro/user/internal/controller/rest/v1"
	"github.com/meraiku/micro/user/internal/service/user"
)

type restService struct {
	userRepo    user.Repository
	userService v1.UserService
	api         *v1.API
	cfg         *config.REST
}

func newRestService() *restService {
	return &restService{}
}

func (s *restService) Run(ctx context.Context) error {

	log := logging.L(ctx)

	log.Debug("initializing rest api")
	api := s.API()

	logging.L(ctx).Info(
		"Starting RestAPI service",
		logging.String("address", s.Config().Address()),
	)

	return api.Run(s.Config().Address())
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

func (s *restService) API() *v1.API {
	if s.api == nil {
		s.api = v1.New(s.Service())
	}

	return s.api
}
