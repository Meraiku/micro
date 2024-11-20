package app

import (
	"context"
	"log"
	"net"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/internal/config"
	"github.com/meraiku/micro/user/internal/containers"
	"github.com/meraiku/micro/user/pkg/auth_v1"
	"github.com/meraiku/micro/user/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type grpcService struct {
	grpcServer *grpc.Server

	authContainer *containers.AuthContainerGRPC
	userContainer *containers.UserContainerGRPC
	cfg           *config.Config
}

func newGRPCService() *grpcService {
	return &grpcService{}
}

func (s *grpcService) Run(ctx context.Context) error {
	log := logging.L(ctx)

	if s.grpcServer == nil {

		log.Debug("initializing grpc server")

		s.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

		reflection.Register(s.grpcServer)

		user_v1.RegisterUserV1Server(s.grpcServer, s.UserContainer().UserAPI)
		auth_v1.RegisterAuthV1Server(s.grpcServer, s.AuthContainer().AuthAPI)
	}

	listner, err := net.Listen("tcp", s.Config().Address())
	if err != nil {
		return err
	}

	log.Info(
		"grpc service initialized",
		logging.String("address", listner.Addr().String()),
	)

	return s.grpcServer.Serve(listner)
}

func (s *grpcService) Config() *config.Config {
	if s.cfg == nil {
		cfg, err := config.NewConfig(context.Background(), config.GRPC)
		if err != nil {
			log.Fatalf("failed to load grpc config: %v", err)
		}

		s.cfg = cfg
	}

	return s.cfg
}

func (s *grpcService) AuthContainer() *containers.AuthContainerGRPC {
	if s.authContainer == nil {
		var err error

		s.authContainer, err = containers.NewAuthGRPC(s.Config())
		if err != nil {
			log.Fatalf("failed to create auth container: %v", err)
		}
	}

	return s.authContainer
}

func (s *grpcService) UserContainer() *containers.UserContainerGRPC {
	if s.userContainer == nil {
		var err error

		s.userContainer, err = containers.NewUserGRPC(s.Config().Repos)
		if err != nil {
			log.Fatalf("failed to create user container: %v", err)
		}
	}

	return s.userContainer
}
