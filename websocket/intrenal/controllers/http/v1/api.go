package v1

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/websocket/intrenal/repo/chatRepo/memory"
	"github.com/meraiku/micro/websocket/intrenal/services/chat"
)

type ChatServiceAPI struct {
	cs   *chat.Service
	addr string
	tmpl *template.Template
}

func NewChatServiceAPI(ctx context.Context, addr string) *ChatServiceAPI {
	tmpl, err := template.New("").ParseFS(templates, "templates/index.html")
	if err != nil {
		log.Fatalf("failed to create template: %v", err)
	}

	repo := memory.NewRepository()
	cs := chat.NewService(ctx, repo)

	return &ChatServiceAPI{
		addr: addr,
		tmpl: tmpl,
		cs:   cs,
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
