package controllers

import (
	"fmt"
	"goforum/admin/helpers"
	amodels "goforum/admin/models"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type Approvals struct {
	Store *sessions.CookieStore
}

// Bekleyen gönderiler listesi
func (a Approvals) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !helpers.CheckUser(w, r, a.Store) {
		return
	}

	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string { return amodels.Category{}.Get(categoryID).Title },
	}).ParseFiles(helpers.Include("approvals/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}

	// admin modelinden bekleyenleri çek (IsApproved=false ve user_id != 0 ise kullanıcı gönderisi)
	pending := amodels.Post{}.GetAll("is_approved = ?", false)

	data := map[string]interface{}{
		"Posts": pending,
		"Alert": helpers.GetAlert(w, r, a.Store),
	}
	view.ExecuteTemplate(w, "index", data)
}

// Bekleyen gönderi detayı
func (a Approvals) Detail(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !helpers.CheckUser(w, r, a.Store) {
		return
	}

	view, err := template.New("index").Funcs(template.FuncMap{
		"safeHTML":    func(html string) template.HTML { return template.HTML(html) },
		"getCategory": func(categoryID int) string { return amodels.Category{}.Get(categoryID).Title },
	}).ParseFiles(helpers.Include("approvals/detail")...)
	if err != nil {
		fmt.Println(err)
		return
	}

	post := amodels.Post{}.Get(p.ByName("id"))
	data := map[string]interface{}{"Post": post, "Alert": helpers.GetAlert(w, r, a.Store)}
	view.ExecuteTemplate(w, "index", data)
}

// Onayla
func (a Approvals) Approve(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !helpers.CheckUser(w, r, a.Store) {
		return
	}
	post := amodels.Post{}.Get(p.ByName("id"))
	if post.ID == 0 {
		helpers.SetAlert(w, r, "Kayıt bulunamadı", a.Store)
		http.Redirect(w, r, "/admin/pending", http.StatusSeeOther)
		return
	}
	post.Update("is_approved", true)
	helpers.SetAlert(w, r, "Gönderi onaylandı", a.Store)
	http.Redirect(w, r, "/admin/pending", http.StatusSeeOther)
}
