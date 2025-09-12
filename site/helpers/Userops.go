package helpers

import (
	adminmodels "goforum/admin/models"
	"net/http"

	"github.com/gorilla/sessions"
)

// SetUser sets the logged-in user's username and id into the cookie session
func SetUser(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore, user adminmodels.User) error {
	session, _ := store.Get(r, "blog-user")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Values["role"] = user.Role
	return session.Save(r, w)
}

// RemoveUser clears the user session
func RemoveUser(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) error {
	session, err := store.Get(r, "blog-user")
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

// GetCurrentUser returns the current user model and a boolean indicating existence
func GetCurrentUser(r *http.Request, store *sessions.CookieStore) (adminmodels.User, bool) {
	session, err := store.Get(r, "blog-user")
	if err != nil {
		return adminmodels.User{}, false
	}
	userID, ok := session.Values["user_id"].(uint)
	if !ok {
		// In some cases Gorilla may encode numeric types differently (e.g., int, int64)
		if id64, ok := session.Values["user_id"].(int64); ok {
			userID = uint(id64)
		} else if idInt, ok := session.Values["user_id"].(int); ok {
			userID = uint(idInt)
		} else {
			return adminmodels.User{}, false
		}
	}
	user := adminmodels.User{Model: adminmodels.User{}.Model}
	user = adminmodels.User{}.Get("id = ?", userID)
	if user.ID == 0 {
		return adminmodels.User{}, false
	}
	return user, true
}
