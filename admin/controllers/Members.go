package controllers

import (
	"fmt"
	"goforum/admin/helpers"
	adminmodels "goforum/admin/models"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type Members struct {
	Store *sessions.CookieStore
}

func (m Members) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !helpers.CheckUser(w, r, m.Store) {
		return
	}

	view, err := template.New("index").Funcs(template.FuncMap{
		"formatTime": func(t time.Time) string { return t.Format("02.01.2006 15:04") },
		"upper":      strings.ToUpper,
	}).ParseFiles(helpers.Include("members/list")...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Tüm kullanıcıları çek
	users := adminmodels.User{}.GetAll()
	var userCount, adminCount int
	for _, u := range users {
		if u.Role == "admin" {
			adminCount++
		} else {
			userCount++
		}
	}

	// Filtre
	filterType := r.URL.Query().Get("type") // "user" | "admin" | ""
	var list []adminmodels.User
	if filterType == "user" {
		for _, u := range users {
			if u.Role != "admin" {
				list = append(list, u)
			}
		}
	} else if filterType == "admin" {
		for _, u := range users {
			if u.Role == "admin" {
				list = append(list, u)
			}
		}
	}

	data := map[string]interface{}{
		"Alert":      helpers.GetAlert(w, r, m.Store),
		"UserCount":  userCount,
		"AdminCount": adminCount,
		"FilterType": filterType,
		"Users":      list,
	}

	_ = view.ExecuteTemplate(w, "index", data)
}
