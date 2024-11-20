package v1

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/pkg/auth_v1"
)

func (s *ChatServiceAPI) authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logging.L(r.Context())

		log.Debug(
			"authenticating user",
		)

		token, err := r.Cookie("access")
		if err != nil {

			log.Debug("no access token found")
			refresh, err := r.Cookie("refresh")
			if err != nil {

				log.Debug("no refresh token found")

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			log.Debug("refreshing tokens")

			tokens, err := s.authSerivce.Refresh(r.Context(), &auth_v1.RefreshRequest{RefreshToken: refresh.Value})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r.AddCookie(&http.Cookie{
				Name:  "access",
				Value: tokens.AccessToken,
			})
			r.AddCookie(&http.Cookie{
				Name:  "refresh",
				Value: tokens.RefreshToken,
			})

			token.Value = tokens.AccessToken
		} else {

			log.Debug("access token found, trying authenticate")

			if _, err := s.authSerivce.Authenticate(r.Context(), &auth_v1.AuthenticateRequest{AccessToken: token.Value}); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (s *ChatServiceAPI) loggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logging.L(r.Context())

		logging.WithAttrs(
			r.Context(),
			logging.String("request_id", uuid.New().String()),
		)

		log.Info(
			"request started",
			logging.String("method", r.Method),
			logging.String("path", r.URL.Path),
		)

		now := time.Now()

		next.ServeHTTP(w, r)

		log.Info(
			"request completed",
			logging.String("method", r.Method),
			logging.String("path", r.URL.Path),
			logging.Duration("duration", time.Since(now)),
		)
	})
}

func (s *ChatServiceAPI) recoverMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log := logging.L(r.Context())

		defer func() {
			if err := recover(); err != nil {
				log.Error(
					"failed to handle request",
					logging.Any("error", err),
				)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *ChatServiceAPI) usernameMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
