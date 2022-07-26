// Server package provides access to HTTP server. It uses chi for routing traffic.
package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server defines dependencies for the HTTP server
type Server struct {
	router chi.Router
}

// New returns a new server
func New() *Server {
	s := &Server{}

	r := chi.NewRouter()

	r.Get("/ping", pong)

	// r.Route("/statistics", func(r chi.Router) {
	// 	h := entryHandler{s.Entry}
	// 	r.Mount("/", h.router())
	// })

	s.router = r
	return s
}

// ServeHTTP represents the http.Handler interface implementation for Server
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, r)
}

func pong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong!\n"))
}

func jsonError(w http.ResponseWriter, details string, code int) {
	err := fmt.Sprintf(`{"error":"%s"}`, details)
	w.WriteHeader(code)
	io.WriteString(w, err)
}
