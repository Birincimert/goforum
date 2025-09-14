package controllers

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	adminhelpers "goforum/site/helpers"
	"goforum/site/models"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type AuthorApplication struct {
	Store *sessions.CookieStore
}

func (aa AuthorApplication) Apply(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(r.FormValue("password"))))
	bio := r.FormValue("bio")

	app := models.AuthorApplication{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Username:  username,
		Password:  password,
		Bio:       bio,
		Status:    "pending",
	}
	app.Add()

	_ = adminhelpers.SetAlert(w, r, "Başvurunuz alındı. Yönetici onayı bekleniyor.", aa.Store)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
