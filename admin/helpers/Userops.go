package helpers

import (
	"goblog/admin/models"
	"net/http"

	"github.com/gorilla/sessions"
)

func SetUser(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore, username string, password string) error {
	session, err := store.Get(r, "blog-user")
	if err != nil {
	}

	session.Values["username"] = username
	session.Values["password"] = password
	return session.Save(r, w)
}

func CheckUser(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) bool {
	session, err := store.Get(r, "blog-user")
	if err != nil {
	}

	username := session.Values["username"]
	password := session.Values["password"]
	user := models.User{}.Get("username = ? AND password = ?", username, password)
	if user.Username == username && user.Password == password {
		return true
	}
	SetAlert(w, r, "Lütfen önce giriş yapın", store)
	returnURL := r.URL.Path
	http.Redirect(w, r, "/admin/login?return_url="+returnURL, http.StatusSeeOther)
	return false
}

func RemoveUser(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) error {
	session, err := store.Get(r, "blog-user")
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
