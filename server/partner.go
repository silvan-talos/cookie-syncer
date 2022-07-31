package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/silvan-talos/cookie-syncer/partner"
)

type partnerHandler struct {
	p partner.Service
}

func (h *partnerHandler) router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.addPartner)
	return r
}

func (h *partnerHandler) addPartner(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var request = struct {
		Name string
		URL  string
	}{
		r.FormValue("Name"),
		r.FormValue("URL"),
	}

	id, err := h.p.AddPartner(request.Name, request.URL)
	if err != nil {
		log.Println("failed to add partner:", err)
		jsonError(w, "failed to add partner", http.StatusInternalServerError)
		return
	}

	var response = struct {
		ID string `json:"partner_uuid"`
	}{
		ID: id,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("failed to encode response:", err)
		jsonError(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
