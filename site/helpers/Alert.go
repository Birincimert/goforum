package helpers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

func SetAlert(w http.ResponseWriter, r *http.Request, message string, store *sessions.CookieStore) error {
	session, err := store.Get(r, "alert-go-site")
	if err != nil {
		fmt.Println(err)
		return err
	}
	session.AddFlash(message)
	return session.Save(r, w)
}

func GetAlert(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) map[string]interface{} {
	session, err := store.Get(r, "alert-go-site")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	data := make(map[string]interface{})
	flashes := session.Flashes()
	if len(flashes) > 0 {
		data["is_alert"] = true
		data["message"] = flashes[0]
	} else {
		data["is_alert"] = false
		data["message"] = nil
	}

	session.Save(r, w)

	return data
}
