package v1

import (
	"embed"
	"net/http"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/pkg/auth_v1"
)

//go:embed templates/*.html
var templates embed.FS

//go:embed css/*.css
var css embed.FS

func (s *ChatServiceAPI) handleRoot(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("HX-Redirect", "/")

	http.ServeFileFS(w, r, templates, "templates/index.html")
}

func (s *ChatServiceAPI) handleUserInfo(w http.ResponseWriter, r *http.Request) {

	http.ServeFileFS(w, r, templates, "templates/user_info.html")
}

func (s *ChatServiceAPI) handleServeCSS(w http.ResponseWriter, r *http.Request) {

	http.ServeFileFS(w, r, css, "templates/css/style.css")
}

func (s *ChatServiceAPI) handleGetChat(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("HX-Redirect", "/chats")

	http.ServeFileFS(w, r, templates, "templates/chat.html")
}

func (s *ChatServiceAPI) handleLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("HX-Redirect", "/login")

	http.ServeFileFS(w, r, templates, "templates/login.html")
}

func (s *ChatServiceAPI) handleRegister(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("HX-Redirect", "/register")

	http.ServeFileFS(w, r, templates, "templates/register.html")
}

func (s *ChatServiceAPI) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	log := logging.L(r.Context())

	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Info(
		"logging in user",
		logging.String("username", username),
	)

	tks, err := s.authSerivce.Login(r.Context(), &auth_v1.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Error(
			"failed to login user",
			logging.Err(err),
		)
		return
	}

	log.Info(
		"user logged in",
		logging.String("username", username),
	)

	log.Debug(
		"setting cookies",
	)

	r.AddCookie(&http.Cookie{
		Name:  "access",
		Value: tks.AccessToken,
	})

	r.AddCookie(&http.Cookie{
		Name:  "refresh",
		Value: tks.RefreshToken,
	})

	log.Debug(
		"redirecting to global chat",
		logging.String("username", username),
	)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *ChatServiceAPI) handleRegisterUser(w http.ResponseWriter, r *http.Request) {

	log := logging.L(r.Context())

	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Info(
		"registering user",
		logging.String("username", username),
	)
	resp, err := s.authSerivce.Register(r.Context(), &auth_v1.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Error(
			"failed to register user",
			logging.Err(err),
		)
		return
	}

	log.Info(
		"user registered",
		logging.String("user_id", resp.Id),
		logging.String("username", resp.Username),
	)

	log.Debug(
		"redirecting to login page",
	)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
