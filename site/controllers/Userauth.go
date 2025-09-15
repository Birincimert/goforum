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
	if user, ok := helpers.GetCurrentUser(r, ua.Store); ok {
		data["CurrentUser"] = user
	}
	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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

	// Add value receiver olduğundan DoRegister içindeki user.ID güncellenmeyebilir.
	// Otomatik giriş için, eklenen kullanıcıyı tekrar çekip (ID dolu) session’a yaz.
	saved := adminmodels.User{}.Get("username = ?", username)
	if saved.ID != 0 {
		_ = helpers.SetUser(w, r, ua.Store, saved)
	} else {
		// Fallback: mevcut user objesi ile (ID 0 olabilir) yine de devam edelim
		_ = helpers.SetUser(w, r, ua.Store, user)
	}
	_ = helpers.SetAlert(w, r, "Kayıt başarılı, hoş geldiniz.", ua.Store)
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
		_ = helpers.SetAlert(w, r, "Hesap bulunamadı, lütfen kayıt olun.", ua.Store)
		http.Redirect(w, r, "/login?return_url="+returnURL, http.StatusSeeOther)
		return
	}
	_ = helpers.SetUser(w, r, ua.Store, user)
	_ = helpers.SetAlert(w, r, "Giriş başarılı.", ua.Store)
	http.Redirect(w, r, returnURL, http.StatusSeeOther)
}

func (ua Userauth) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_ = helpers.RemoveUser(w, r, ua.Store)
	_ = helpers.SetAlert(w, r, "Çıkış yapıldı.", ua.Store)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
