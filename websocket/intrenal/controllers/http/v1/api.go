package v1

import (
	"context"
	"net/http"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/websocket/intrenal/models"
	"github.com/meraiku/micro/websocket/intrenal/services/chat"
)

type AuthService interface {
	Login(ctx context.Context, user *models.User) (*models.Tokens, error)
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Authenticate(ctx context.Context, accessToken string) (*models.User, error)
	Refresh(ctx context.Context, refreshToken string) (*models.Tokens, error)
}

type ChatService interface {
	ConnectGlobal(client *chat.Client, w http.ResponseWriter, r *http.Request) error
}

type ChatServiceAPI struct {
	cs          ChatService
	authService AuthService

	addr string
}

func NewChatServiceAPI(
	ctx context.Context,
	addr string,
	chatService ChatService,
	authService AuthService,
) *ChatServiceAPI {

	return &ChatServiceAPI{
		addr:        addr,
		authService: authService,
		cs:          chatService,
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
