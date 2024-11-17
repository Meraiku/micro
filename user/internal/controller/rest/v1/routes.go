package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *API) routes() *chi.Mux {

	r := chi.NewRouter()

	r.Use(loggingMiddleware)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.Get("/users", api.Make(api.GetUsers))
			r.Post("/users", api.Make(api.CreateUser))
			r.Get("/users/{id}", api.Make(api.GetUserByID))
			r.Put("/users/{id}", api.Make(api.UpdateUser))
			r.Delete("/users/{id}", api.Make(api.DeleteUser))
		})
	})

	return r
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func (api *API) Make(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {

			if apiErr, ok := err.(APIError); ok {
				api.JSON(w, apiErr.StatusCode, apiErr)
			} else {
				errResp := map[string]any{
					"status_code": http.StatusInternalServerError,
					"msg":         "internal server error",
				}
				api.JSON(w, http.StatusInternalServerError, errResp)
			}
		}
	}
}

type APIError struct {
	StatusCode int `json:"status_code" example:"400"`
	Msg        any `json:"msg" swaggertype:"string" example:"invalid ID"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %s", e.Msg)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func InvalidRequestData(errors map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        errors,
	}
}
