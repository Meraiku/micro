package app

import (
	"context"
	"log"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/internal/config"
	"github.com/meraiku/micro/user/internal/containers"
	v1 "github.com/meraiku/micro/user/internal/controller/rest/v1"
)

type restService struct {
	api *v1.API
	cfg *config.Config

	userContainer *containers.UserContainerREST
}

func newRestService() *restService {
	return &restService{}
}

func (s *restService) Run(ctx context.Context) error {

	log := logging.L(ctx)

	log.Debug("initializing rest api")

	api := s.UserContainer().UserAPI

	log.Info(
		"rest service initialized",
		logging.String("address", s.Config().Address()),
	)

	return api.Run(s.Config().Address())
}

func (s *restService) Config() *config.Config {
	if s.cfg == nil {
		cfg, err := config.NewConfig(context.Background(), config.REST)
		if err != nil {
			log.Fatalf("failed to load rest config: %v", err)
		}

		s.cfg = cfg
	}

	return s.cfg
}

func (s *restService) UserContainer() *containers.UserContainerREST {
	if s.userContainer == nil {
		var err error

		s.userContainer, err = containers.NewUserREST(s.Config().Repos)
		if err != nil {
			log.Fatalf("failed to create user container: %v", err)
		}
	}

	return s.userContainer
}
