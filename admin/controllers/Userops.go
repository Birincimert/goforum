package controllers

import (
	"crypto/sha256"
	"fmt"
	"goblog/admin/helpers"
	"goblog/admin/models"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type Userops struct {
	Store *sessions.CookieStore
}

func (userops Userops) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	view, err := template.ParseFiles(helpers.Include("userops/login")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Alert"] = helpers.GetAlert(w, r, userops.Store)
	data["ReturnURL"] = r.URL.Query().Get("return_url")
	view.ExecuteTemplate(w, "index", data)
}

func (userops Userops) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	username := r.FormValue("username")
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(r.FormValue("password"))))
	returnURL := r.FormValue("return_url")

	if returnURL == "" {
		returnURL = "/admin"
	}

	user := models.User{}.Get("username = ? AND password = ?", username, password)
	if user.Username == username && user.Password == password {
		// Rol kontrolü ekle: Sadece 'admin' rolüne sahip kullanıcılar giriş yapabilsin
		if user.Role == "admin" {
			helpers.SetUser(w, r, userops.Store, username, password)
			helpers.SetAlert(w, r, "Hoşgeldiniz.", userops.Store)
			http.Redirect(w, r, returnURL, http.StatusSeeOther)
		} else {
			// Rolü admin değilse hata mesajı ver
			helpers.SetAlert(w, r, "Bu alana erişim yetkiniz yoktur.", userops.Store)
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		}
	} else {
		helpers.SetAlert(w, r, "Yanlış kullanıcı adı veya şifre.", userops.Store)
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

}

func (userops Userops) Logout(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	helpers.RemoveUser(w, r, userops.Store)
	helpers.SetAlert(w, r, "Hoşçakalın!", userops.Store)
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
