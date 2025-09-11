package controllers

import (
	"goblog/admin/helpers"
	"goblog/site/models"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type Comments struct {
	Store *sessions.CookieStore
}

// Index: List posts with their comment counts
func (comments Comments) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Yazıları yorumlarıyla birlikte al
	posts := models.Post{}.GetAll()

	// Görünüm verisi hazırlanır: her yazı için toplam yorum sayısı (yanıtlar dahil)
	type postWithCount struct {
		ID           uint
		Title        string
		Slug         string
		CommentCount int
	}

	var list []postWithCount
	for _, p := range posts {
		total := 0
		for _, c := range p.Comments {
			total++
			total += len(c.Replies)
		}
		list = append(list, postWithCount{ID: p.ID, Title: p.Title, Slug: p.Slug, CommentCount: total})
	}

	view := make(map[string]interface{})
	view["Posts"] = list
	view["Alert"] = helpers.GetAlert(w, r, comments.Store)

	files := helpers.Include("comments/posts")
	temp := template.Must(template.ParseFiles(files...))
	if err := temp.ExecuteTemplate(w, "index", view); err != nil {
		http.Error(w, "Template hatası: "+err.Error(), http.StatusInternalServerError)
	}
}

func (comments Comments) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ID parametresini al ve dönüştür
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		http.Error(w, "Geçersiz yorum ID", http.StatusBadRequest)
		return
	}

	// Yorumu sil
	comment := models.Comment{}
	if err := comment.Delete(uint(id)); err != nil {
		http.Error(w, "Yorum silinirken hata oluştu", http.StatusInternalServerError)
		return
	}

	// Alert ayarla ve geri yönlendir
	_ = helpers.SetAlert(w, r, "Yorum başarıyla silindi", comments.Store)
	back := r.Header.Get("Referer")
	if back == "" {
		back = "/admin/comments"
	}
	http.Redirect(w, r, back, http.StatusSeeOther)
}

// Post: List comments for a specific post
func (comments Comments) Post(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Geçersiz yazı ID", http.StatusBadRequest)
		return
	}

	// Yazıyı ve yorumlarını al
	post := models.Post{}.Get("id = ?", uint(id))
	if post.ID == 0 {
		http.NotFound(w, r)
		return
	}

	// Sadece görünüme uygun basit alanları verelim
	type viewPost struct {
		ID    uint
		Title string
		Slug  string
	}

	view := make(map[string]interface{})
	view["Post"] = viewPost{ID: post.ID, Title: post.Title, Slug: post.Slug}
	view["Comments"] = post.Comments
	view["Alert"] = helpers.GetAlert(w, r, comments.Store)

	files := helpers.Include("comments/post")
	temp := template.Must(template.ParseFiles(files...))
	if err := temp.ExecuteTemplate(w, "index", view); err != nil {
		http.Error(w, "Template hatası: "+err.Error(), http.StatusInternalServerError)
	}
}
