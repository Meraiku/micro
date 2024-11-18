package v1

import (
	"embed"
	"log"
	"net/http"

	"github.com/meraiku/micro/user/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed templates/*.html
var templates embed.FS

func (s *ChatServiceAPI) handleGetChat(w http.ResponseWriter, r *http.Request) {

	http.ServeFileFS(w, r, templates, "templates/index.html")
}

func (s *ChatServiceAPI) handleLogin(w http.ResponseWriter, r *http.Request) {

	http.ServeFileFS(w, r, templates, "templates/login.html")
}

func (s *ChatServiceAPI) handleLoginUser(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func (s *ChatServiceAPI) handleRegisterUser(w http.ResponseWriter, r *http.Request) {

	conn, err := grpc.NewClient("user_service_grpc:20001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create grpc client: %v", err)
		return
	}

	client := auth_v1.NewAuthV1Client(conn)

	resp, err := client.Register(r.Context(), &auth_v1.RegisterRequest{
		Username: r.FormValue("register-username"),
		Password: r.FormValue("register-password"),
	})
	if err != nil {
		log.Printf("failed to register user: %v", err)
		return
	}

	log.Printf("user registered: %v", resp)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
