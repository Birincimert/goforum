package controllers

import (
	"fmt"
	"goblog/admin/helpers"
	"goblog/admin/models"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/gosimple/slug"
	"github.com/julienschmidt/httprouter"
)

type Categories struct {
	Store *sessions.CookieStore
}

func (categories Categories) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !helpers.CheckUser(w, r, categories.Store) {
		return
	}

	view, err := template.ParseFiles(helpers.Include("categories/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := map[string]interface{}{}
	data["Categories"] = models.Category{}.GetAll()
	data["Alert"] = helpers.GetAlert(w, r, categories.Store)
	view.ExecuteTemplate(w, "index", data)
}

func (categories Categories) Add(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.CheckUser(w, r, categories.Store) {
		return
	}

	categoryTitle := r.FormValue("category-title")
	categorySlug := slug.Make(categoryTitle)

	models.Category{Title: categoryTitle, Slug: categorySlug}.Add()
	helpers.SetAlert(w, r, "Kategori başarıyla eklendi.", categories.Store)
	http.Redirect(w, r, "/admin/kategoriler", http.StatusSeeOther)
}

func (categories Categories) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.CheckUser(w, r, categories.Store) {
		return
	}

	category := models.Category{}.Get(params.ByName("id"))
	category.Delete()

	helpers.SetAlert(w, r, "Kategori başarıyla silindi...", categories.Store)
	http.Redirect(w, r, "/admin/kategoriler", http.StatusSeeOther)
}
