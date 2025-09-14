package controllers

import (
	"fmt"
	tmpl "html/template"
	"net/http"
	"strconv"

	"goforum/admin/helpers"
	adminmodels "goforum/admin/models"
	sitemodels "goforum/site/models"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type Members struct {
	Store *sessions.CookieStore
}

func (m Members) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view, err := tmpl.ParseFiles(helpers.Include("members/list")...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := make(map[string]interface{})
	data["Alert"] = helpers.GetAlert(w, r, m.Store)

	users := adminmodels.User{}.GetAll()
	// only pending applications
	apps := sitemodels.AuthorApplication{}.GetAll("status = ?", "pending")

	// Split users by role
	var admins []adminmodels.User
	var authors []adminmodels.User
	var regulars []adminmodels.User
	for _, u := range users {
		switch u.Role {
		case "admin":
			admins = append(admins, u)
		case "author":
			authors = append(authors, u)
		default:
			regulars = append(regulars, u)
		}
	}

	data["Admins"] = admins
	data["Yazars"] = authors
	data["Users"] = regulars
	data["Applications"] = apps

	_ = view.ExecuteTemplate(w, "index", data)
}

func (m Members) Approve(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	idStr := params.ByName("id")
	id, _ := strconv.Atoi(idStr)
	app := sitemodels.AuthorApplication{}.Get("id = ?", id)
	if app.ID == 0 {
		_ = helpers.SetAlert(w, r, "Başvuru bulunamadı.", m.Store)
		http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
		return
	}

	// Create user as author
	user := adminmodels.User{
		FirstName: app.FirstName,
		LastName:  app.LastName,
		Email:     app.Email,
		Username:  app.Username,
		Password:  app.Password,
		Role:      "author",
	}
	user.Add()
	// Update application status
	app.UpdateStatus("approved")

	helperMsg := "Başvuru onaylandı, kullanıcı yazar olarak kaydedildi."
	_ = helpers.SetAlert(w, r, helperMsg, m.Store)
	http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
}

func (m Members) Reject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	idStr := params.ByName("id")
	id, _ := strconv.Atoi(idStr)
	app := sitemodels.AuthorApplication{}.Get("id = ?", id)
	if app.ID == 0 {
		_ = helpers.SetAlert(w, r, "Başvuru bulunamadı.", m.Store)
		http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
		return
	}
	app.UpdateStatus("rejected")
	_ = helpers.SetAlert(w, r, "Başvuru reddedildi.", m.Store)
	http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
}
