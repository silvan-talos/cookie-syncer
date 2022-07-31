package server

import (
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
	log.Println("Start syncing")
	otherPartner := make(chan syncer.Partner)
	go h.getOtherPartner(otherPartner, partnerID)
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
		}
	} else {
		userID = c.Value
		log.Println("cookie already there, userID:", userID)
	}
	c.MaxAge = 3600
	http.SetCookie(w, c)
	op := <-otherPartner
	otherPID := op.ID
	opURL := op.URL
	url := opURL + "/partners/" + otherPID + "/sync/img.gif?partner=" + partnerID + "&partner_uid=" + userID
	log.Println("url:", url, "opurl", opURL)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}
func (h *syncingHandler) saveMatchAndRedirect(partnerID, otherPID, opUserID string, w http.ResponseWriter, r *http.Request) {
	log.Println("Save and redirect")
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
		}
	} else {
		userID = c.Value
		log.Println("cookie already there, userID:", userID)
	}
	//save Match
	go h.s.SyncUsers(partnerID, otherPID, userID, opUserID)
	c.MaxAge = 3600
	http.SetCookie(w, c)
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
	c.MaxAge = 3600
	http.SetCookie(w, c)
	w.WriteHeader(http.StatusNoContent)
}
