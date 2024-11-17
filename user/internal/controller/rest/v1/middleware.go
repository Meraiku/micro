package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/meraiku/micro/pkg/logging"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := uuid.New()

		log := logging.L(r.Context()).With("request_id", requestID)

		ctx := logging.ContextWithLogger(r.Context(), log)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		log.Info(
			"request completed",
			logging.String("method", r.Method),
			logging.String("path", r.URL.Path),
		)
	})
}
