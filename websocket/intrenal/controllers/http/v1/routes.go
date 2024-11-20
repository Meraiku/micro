package v1

import (
	"github.com/go-chi/chi/v5"
)

func (s *ChatServiceAPI) setupRoutes() *chi.Mux {

	r := chi.NewRouter()

	r.Use(s.loggingMiddleware)
	r.Use(s.recoverMiddleware)

	r.Get("/", s.handleRoot)
	FileServerFS(r, "/", css)

	r.Get("/login", s.handleLogin)
	r.Post("/login", s.handleLoginUser)

	r.Get("/register", s.handleRegister)
	r.Post("/register", s.handleRegisterUser)

	r.Get("/users", s.handleUserInfo)

	r.Route("/chats", func(r chi.Router) {

		r.Use(s.authMiddleware)

		r.Get("/", s.handleGetChat)
		r.Get("/global", s.handleGlobalChat)
	})

	return r
}
