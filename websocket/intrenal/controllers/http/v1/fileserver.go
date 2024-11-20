package v1

import (
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func FileServerFS(r chi.Router, path string, root fs.FS) {
	if strings.ContainsAny(path, "{}*") {
		log.Printf("FileServer does not permit URL parameters.")
		return
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		prefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")

		fs := http.StripPrefix(prefix, http.FileServerFS(root))
		fs.ServeHTTP(w, r)
	})
}
