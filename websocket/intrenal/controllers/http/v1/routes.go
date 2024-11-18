package v1

import "github.com/go-chi/chi/v5"

func (s *ChatServiceAPI) setupRoutes() *chi.Mux {

	r := chi.NewRouter()

	r.Get("/login", s.handleLogin)

	r.Route("/chat", func(r chi.Router) {
		r.Get("/", s.handleGetChat)
		r.Get("/ws", s.handleWebsocket)
	})

	return r
}
