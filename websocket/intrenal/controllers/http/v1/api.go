package v1

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/pkg/auth_v1"
	"github.com/meraiku/micro/websocket/intrenal/repo/chatRepo/memory"
	"github.com/meraiku/micro/websocket/intrenal/services/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatServiceAPI struct {
	cs          *chat.Service
	addr        string
	authAddr    string
	authSerivce auth_v1.AuthV1Client
}

func NewChatServiceAPI(ctx context.Context, addr string) *ChatServiceAPI {

	repo := memory.NewRepository()
	cs := chat.NewService(ctx, repo)

	authAddr := os.Getenv("AUTH_SERVICE")
	if authAddr == "" {
		authAddr = "http://localhost:20001"
	}

	conn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create grpc client: %v", err)
	}

	authSerivce := auth_v1.NewAuthV1Client(conn)

	return &ChatServiceAPI{
		addr:        addr,
		authAddr:    authAddr,
		authSerivce: authSerivce,
		cs:          cs,
	}
}

func (s *ChatServiceAPI) Run(ctx context.Context) error {

	srv := http.Server{
		Addr:    s.addr,
		Handler: s.setupRoutes(),
	}

	logging.L(ctx).Info(
		"Running HTTP server",
		logging.String("address", s.addr),
	)

	return srv.ListenAndServe()
}
