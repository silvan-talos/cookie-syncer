package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	syncer "github.com/silvan-talos/cookie-syncer"
	"github.com/silvan-talos/cookie-syncer/partner"
	"github.com/silvan-talos/cookie-syncer/syncing"
)

type syncingHandler struct {
	p partner.Service
	s syncing.Service
}

func (h *syncingHandler) router() chi.Router {
	r := chi.NewRouter()
	r.Get("/img.gif", h.getPixel)
	r.Get("/usermatch.gif", h.saveUserMatch)
	r.Get("/", h.getSyncStatus)
	return r
}

func (h *syncingHandler) getPixel(w http.ResponseWriter, r *http.Request) {
	partnerID := chi.URLParam(r, "partnerID")
	log.Println("getPixel for partnerID:", partnerID)
	otherPID := r.URL.Query().Get("partner")
	opUserID := r.URL.Query().Get("partner_uid")
	if otherPID == "" || opUserID == "" {
		h.syncWithOther(partnerID, w, r)
		return
	}
	h.saveMatchAndRedirect(partnerID, otherPID, opUserID, w, r)
}

func (h *syncingHandler) syncWithOther(partnerID string, w http.ResponseWriter, r *http.Request) {
	otherPartner := make(chan syncer.Partner)
	go h.getOtherPartner(otherPartner, partnerID)
	userID := getUserIDfromCookie(partnerID, w, r)
	op := <-otherPartner
	otherPID := op.ID
	opURL := op.URL
	url := opURL + "/partners/" + otherPID + "/sync/img.gif?partner=" + partnerID + "&partner_uid=" + userID
	log.Println("url:", url)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}

func getUserIDfromCookie(partnerID string, w http.ResponseWriter, r *http.Request) string {
	userID := ""
	c, err := r.Cookie(partnerID)
	if err != nil {
		log.Println("Cookie not found. Adding one...")
		userID = uuid.New().String()
		c = &http.Cookie{
			Name:     partnerID,
			Value:    userID,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			Path:     "/",
			MaxAge:   3600 * 24 * 365,
		}
	} else {
		userID = c.Value
		log.Println("cookie already there, userID:", userID)
	}
	http.SetCookie(w, c)
	return userID
}

func (h *syncingHandler) saveMatchAndRedirect(partnerID, otherPID, opUserID string, w http.ResponseWriter, r *http.Request) {
	log.Println("Save and redirect")
	userID := getUserIDfromCookie(partnerID, w, r)
	//save Match
	go h.s.SyncUsers(partnerID, otherPID, userID, opUserID)
	otherPartner, err := h.p.GetByID(otherPID)
	if err != nil {
		log.Println("failed to get partner with id", otherPID, "error:", err)
		jsonError(w, "failed to get partner", http.StatusBadRequest)
		return
	}
	url := otherPartner.URL + "/partners/" + otherPID + "/sync/usermatch.gif?partner=" + partnerID + "&partner_uid=" + userID
	log.Println("url:", url)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}

func (h *syncingHandler) getOtherPartner(p chan<- syncer.Partner, partnerID string) {
	partners := h.p.GetAll()
	for _, partner := range partners {
		if partner.ID != partnerID {
			p <- partner
			break
		}
	}
}

func (h *syncingHandler) saveUserMatch(w http.ResponseWriter, r *http.Request) {
	log.Println("save callback")
	partnerID := chi.URLParam(r, "partnerID")
	otherPID := r.URL.Query().Get("partner")
	opUserID := r.URL.Query().Get("partner_uid")
	if otherPID == "" || opUserID == "" {
		log.Println("query parameters missing")
		jsonError(w, "query parameters missing", http.StatusBadRequest)
		return
	}
	userID := ""
	c, err := r.Cookie(partnerID)
	if err != nil {
		log.Println("Cookie not found. Aborting...")
		jsonError(w, "cookie not found", http.StatusBadRequest)
		return
	} else {
		userID = c.Value
		log.Println("cookie is there, userID:", userID)
	}
	//save Match
	go h.s.SyncUsers(partnerID, otherPID, userID, opUserID)
	http.SetCookie(w, c)
	w.WriteHeader(http.StatusNoContent)
}

func (h *syncingHandler) getSyncStatus(w http.ResponseWriter, r *http.Request) {
	partnerID := chi.URLParam(r, "partnerID")
	syncs, err := h.s.GetStatus(partnerID)
	if err != nil {
		log.Println("failed to get status from db:", err)
		jsonError(w, "no data available", http.StatusNotFound)
		return
	}
	var response = struct {
		Syncs []syncing.Status `json:"status"`
	}{
		Syncs: syncs,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("error encoding status response:", err)
		jsonError(w, "error getting status", http.StatusInternalServerError)
		return
	}
}
