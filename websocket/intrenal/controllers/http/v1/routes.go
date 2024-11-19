package v1

import "github.com/go-chi/chi/v5"

func (s *ChatServiceAPI) setupRoutes() *chi.Mux {

	r := chi.NewRouter()

	r.Get("/login", s.handleLogin)
	r.Post("/login", s.handleLoginUser)

	r.Post("/register", s.handleRegisterUser)

	r.Route("/chat", func(r chi.Router) {
		r.Get("/", s.handleGetChat)
		r.Get("/global", s.handleGlobalChat)
	})

	return r
}
