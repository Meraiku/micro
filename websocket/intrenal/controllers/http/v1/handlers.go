package v1

import (
	"embed"
	"net/http"
)

//go:embed templates/*.html
var templates embed.FS

func (s *ChatServiceAPI) handleGetChat(w http.ResponseWriter, r *http.Request) {

	http.ServeFileFS(w, r, templates, "templates/index.html")
}

func (s *ChatServiceAPI) handleLogin(w http.ResponseWriter, r *http.Request) {

	http.ServeFileFS(w, r, templates, "templates/login.html")
}
