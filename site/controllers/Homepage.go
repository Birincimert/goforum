package controllers

import (
	"fmt"
	"goforum/site/helpers"
	"goforum/site/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type Homepage struct {
	Store *sessions.CookieStore
}

func (homepage Homepage) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string {
			return models.Category{}.Get(categoryID).Title
		},
		"getDate": func(t time.Time) string {
			return fmt.Sprintf("%02d.%02d.%d", t.Day(), int(t.Month()), t.Year())
		},
		"sumReplies": func(comments []models.Comment) int {
			var total int
			var walk func(items []models.Comment)
			walk = func(items []models.Comment) {
				for _, c := range items {
					total++
					if len(c.Replies) > 0 {
						walk(c.Replies)
					}
				}
			}
			walk(comments)
			return total
		},
	}).ParseFiles(helpers.Include("homepage/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Posts"] = models.Post{}.GetAll()
	// include alert so SetAlert redirects show toast on homepage
	data["Alert"] = helpers.GetAlert(w, r, homepage.Store)
	if user, ok := helpers.GetCurrentUser(r, homepage.Store); ok {
		data["CurrentUser"] = user
	}
	data["ReturnURL"] = r.URL.RequestURI()
	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (homepage Homepage) Detail(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	funcMap := template.FuncMap{
		"safeHTML": func(html string) template.HTML {
			return template.HTML(html)
		},
		"sumReplies": func(comments []models.Comment) int {
			var total int
			var walk func(items []models.Comment)
			walk = func(items []models.Comment) {
				for _, c := range items {
					total++
					if len(c.Replies) > 0 {
						walk(c.Replies)
					}
				}
			}
			walk(comments)
			return total
		},
		"add": func(a int, b int) int { return a + b },
	}
	view, err := template.New("detail").Funcs(funcMap).ParseFiles(helpers.Include("homepage/detail")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	slug := params.ByName("slug")
	post := models.Post{}.GetBySlug(slug)
	if post.ID == 0 {
		if id, errNum := strconv.ParseUint(slug, 10, 64); errNum == nil {
			p := models.Post{}.Get(id)
			if p.ID != 0 {
				post = p // render with ID fallback if slug missing
			}
		}
	}
	data["Post"] = post
	if user, ok := helpers.GetCurrentUser(r, homepage.Store); ok {
		data["CurrentUser"] = user
	}
	data["Alert"] = helpers.GetAlert(w, r, homepage.Store)
	// Use canonical slug if available for return URL
	returnSlug := slug
	if post.ID != 0 && post.Slug != "" {
		returnSlug = post.Slug
	}
	data["ReturnURL"] = "/post/" + returnSlug + "#comments"
	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (homepage Homepage) About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Veritabanından veriyi çek
	about := models.About{}.Get()

	// 2. safeHTML fonksiyonunu şablon motoruna ekle
	funcMap := template.FuncMap{
		"safeHTML": func(html string) template.HTML {
			return template.HTML(html)
		},
	}

	// 3. Şablon dosyasını ana şablonla birlikte işle
	view, err := template.New("about").Funcs(funcMap).ParseFiles(
		helpers.Include("about")...,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 4. Veriyi şablona gönder
	data := map[string]interface{}{}
	data["About"] = about
	if user, ok := helpers.GetCurrentUser(r, homepage.Store); ok {
		data["CurrentUser"] = user
	}
	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Profile shows a simple profile page for logged-in users
func (homepage Homepage) Profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, ok := helpers.GetCurrentUser(r, homepage.Store)
	if !ok {
		http.Redirect(w, r, "/login?return_url=/profile", http.StatusSeeOther)
		return
	}
	view, err := template.ParseFiles(helpers.Include("homepage/profile")...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"CurrentUser": user,
	}
	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (homepage Homepage) Contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view, err := template.ParseFiles(helpers.Include("/contact")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := map[string]interface{}{}
	if user, ok := helpers.GetCurrentUser(r, homepage.Store); ok {
		data["CurrentUser"] = user
	}
	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
