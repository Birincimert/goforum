package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"goforum/site/helpers"
	sitemodels "goforum/site/models"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type SaveController struct {
	Store *sessions.CookieStore
}

// ToggleSave: yazıyı kaydet/kaldır
func (c SaveController) ToggleSave(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, ok := helpers.GetCurrentUser(r, c.Store)
	if !ok {
		http.Error(w, "Giriş gerekli", http.StatusUnauthorized)
		return
	}
	idStr := ps.ByName("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Geçersiz post ID", http.StatusBadRequest)
		return
	}
	saved, err := (sitemodels.PostSave{}).Toggle(user.ID, uint(id64))
	if err != nil {
		http.Error(w, "İşlem başarısız", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"success": true, "saved": %t}`, saved)
}

// IsSaved: kullanıcı bu yazıyı kaydetmiş mi?
func (c SaveController) IsSaved(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, ok := helpers.GetCurrentUser(r, c.Store)
	if !ok {
		http.Error(w, "Giriş gerekli", http.StatusUnauthorized)
		return
	}
	idStr := ps.ByName("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Geçersiz post ID", http.StatusBadRequest)
		return
	}
	isSaved, err := (sitemodels.PostSave{}).IsSaved(user.ID, uint(id64))
	if err != nil {
		http.Error(w, "Durum alınamadı", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"success": true, "saved": %t}`, isSaved)
}
