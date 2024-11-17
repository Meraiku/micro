package app

import (
	"log"
	"net"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/internal/config"
	v1 "github.com/meraiku/micro/user/internal/controller/grpc/v1"
	"github.com/meraiku/micro/user/internal/service/user"
	"github.com/meraiku/micro/user/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type grpcService struct {
	grpcServer *grpc.Server

	userRepo    user.Repository
	userService v1.UserService
	api         *v1.GRPCServer
	cfg         *config.GRPC
}

func newGRPCService() *grpcService {
	return &grpcService{}
}

func (s *grpcService) Run() error {
	if s.grpcServer == nil {
		s.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

		reflection.Register(s.grpcServer)

		user_v1.RegisterUserV1Server(s.grpcServer, s.API())
	}

	listner, err := net.Listen("tcp", s.Config().Address())
	if err != nil {
		return err
	}

	logging.Default().Info(
		"grpc service initialized",
		logging.StringAttr("address", listner.Addr().String()),
	)

	return s.grpcServer.Serve(listner)
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

func (s *grpcService) API() *v1.GRPCServer {
	if s.api == nil {
		s.api = v1.New(s.Service())
	}

	return s.api
}
