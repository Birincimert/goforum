package controllers

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"net/http"

	adminmodels "goforum/admin/models"
	"goforum/site/helpers"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type Userauth struct {
	Store *sessions.CookieStore
}

func (ua Userauth) LoginRegisterPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view, err := template.ParseFiles(helpers.Include("userauth")...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{}
	data["Alert"] = helpers.GetAlert(w, r, ua.Store)
	data["ReturnURL"] = r.URL.Query().Get("return_url")
	view.ExecuteTemplate(w, "index", data)
}

func (ua Userauth) DoRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(r.FormValue("password"))))
	returnURL := r.FormValue("return_url")
	if returnURL == "" {
		returnURL = "/"
	}

	user := adminmodels.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Username:  username,
		Password:  password,
		Role:      "user",
	}
	user.Add()

	helpers.SetUser(w, r, ua.Store, user)
	helpers.SetAlert(w, r, "Kayıt başarılı, hoş geldiniz.", ua.Store)
	http.Redirect(w, r, returnURL, http.StatusSeeOther)
}

func (ua Userauth) DoLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(r.FormValue("password"))))
	returnURL := r.FormValue("return_url")
	if returnURL == "" {
		returnURL = "/"
	}
	user := adminmodels.User{}.Get("username = ? AND password = ?", username, password)
	if user.ID == 0 {
		helpers.SetAlert(w, r, "Hesap bulunamadı, lütfen kayıt olun.", ua.Store)
		http.Redirect(w, r, "/login?return_url="+returnURL, http.StatusSeeOther)
		return
	}
	helpers.SetUser(w, r, ua.Store, user)
	helpers.SetAlert(w, r, "Giriş başarılı.", ua.Store)
	http.Redirect(w, r, returnURL, http.StatusSeeOther)
}

func (ua Userauth) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	helpers.RemoveUser(w, r, ua.Store)
	helpers.SetAlert(w, r, "Çıkış yapıldı.", ua.Store)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
