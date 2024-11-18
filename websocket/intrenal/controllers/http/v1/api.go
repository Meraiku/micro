package v1

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/meraiku/micro/pkg/logging"
)

type ChatServiceAPI struct {
	addr string
	tmpl *template.Template
	hub  *Hub
}

func NewChatServiceAPI(addr string) *ChatServiceAPI {
	tmpl, err := template.New("").ParseFS(templates, "templates/index.html")
	if err != nil {
		log.Fatalf("failed to create template: %v", err)
	}

	h := NewHub()

	go h.run()

	return &ChatServiceAPI{
		addr: addr,
		tmpl: tmpl,
		hub:  h,
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
