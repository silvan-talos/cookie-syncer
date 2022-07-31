// Server package provides access to HTTP server. It uses chi for routing traffic.
package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/silvan-talos/cookie-syncer/partner"
	"github.com/silvan-talos/cookie-syncer/syncing"
)

// Server defines dependencies for the HTTP server
type Server struct {
	Partner partner.Service
	Syncing syncing.Service

	router chi.Router
}

// New returns a new server
func New(ps partner.Service, ss syncing.Service) *Server {
	s := &Server{
		Partner: ps,
		Syncing: ss,
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/ping", pong)

	r.Route("/partners", func(r chi.Router) {
		h := partnerHandler{s.Partner}
		r.Mount("/", h.router())
	})

	r.Route("/partners/{partnerID}/sync", func(r chi.Router) {
		h := syncingHandler{s.Partner, s.Syncing}
		r.Mount("/", h.router())
	})

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
