package controllers

import (
	"fmt"
	adminmodels "goforum/admin/models"
	"goforum/site/helpers"
	"goforum/site/models"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/gosimple/slug"
	"github.com/julienschmidt/httprouter"
)

type Homepage struct {
	Store *sessions.CookieStore
}

func (homepage Homepage) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Kategori sayıları
	counts := models.Post{}.CountsByCategory()

	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string { return models.Category{}.Get(categoryID).Title },
		"getDate":     func(t time.Time) string { return fmt.Sprintf("%02d.%02d.%d", t.Day(), int(t.Month()), t.Year()) },
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
		// Yeni: yazarı getir ve admin post kontrolü
		"getAuthor": func(userID uint) adminmodels.User {
			if userID == 0 {
				return adminmodels.User{}
			}
			return adminmodels.User{}.Get("id = ?", userID)
		},
		"isAdminPost": func(userID uint) bool { return userID == 0 },
		"firstLetter": func(s string) string {
			if s == "" {
				return "?"
			}
			r := []rune(s)
			return strings.ToUpper(string(r[0]))
		},
		// Kategori yazı sayısı yardımcıları
		"catCount": func(catID uint) int { return counts[int(catID)] },
		"totalCount": func() int {
			sum := 0
			for _, v := range counts {
				sum += v
			}
			return sum
		},
	}).ParseFiles(helpers.Include("homepage/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := make(map[string]interface{})
	// SADECE onaylı veya admin yazıları göster
	data["Posts"] = models.Post{}.GetAll("is_approved = ? OR user_id = 0", true)
	// Kategorileri şablona gönder
	data["Categories"] = models.Category{}.GetAll()
	// Aktif kategori (tümü)
	data["ActiveCategorySlug"] = ""
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
		// Detay şablonunda kullanılan yardımcılar
		"getAuthor": func(userID uint) adminmodels.User {
			if userID == 0 {
				return adminmodels.User{}
			}
			return adminmodels.User{}.Get("id = ?", userID)
		},
		"isAdminPost": func(userID uint) bool { return userID == 0 },
		"firstLetter": func(s string) string {
			if s == "" {
				return "?"
			}
			r := []rune(s)
			return strings.ToUpper(string(r[0]))
		},
		"getDate": func(t time.Time) string {
			return fmt.Sprintf("%02d.%02d.%d", t.Day(), int(t.Month()), t.Year())
		},
	}
	view, err := template.New("detail").Funcs(funcMap).ParseFiles(helpers.Include("homepage/detail")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	slugStr := params.ByName("slug")
	post := models.Post{}.GetBySlug(slugStr)
	if post.ID == 0 {
		if id, errNum := strconv.ParseUint(slugStr, 10, 64); errNum == nil {
			p := models.Post{}.Get(id)
			if p.ID != 0 {
				post = p // render with ID fallback if slug missing
			}
		}
	}
	// Onaylanmamış kullanıcı gönderilerini herkese açık göstermeyelim
	if post.ID != 0 && post.UserID != 0 && !post.IsApproved {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	data["Post"] = post
	if user, ok := helpers.GetCurrentUser(r, homepage.Store); ok {
		data["CurrentUser"] = user
	}
	data["Alert"] = helpers.GetAlert(w, r, homepage.Store)
	// Use canonical slug if available for return URL
	returnSlug := slugStr
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
	// Kullanıcının yazıları
	myPosts := models.Post{}.GetAll("user_id = ?", user.ID)
	// Kullanıcının yorumları
	myComments, _ := models.Comment{}.GetByUserID(user.ID)
	// Beğendiği yorumlar
	likedComments, _ := models.GetLikedCommentsByUser(user.ID)

	data := map[string]interface{}{
		"CurrentUser":   user,
		"MyPosts":       myPosts,
		"MyComments":    myComments,
		"LikedComments": likedComments,
		"Alert":         helpers.GetAlert(w, r, homepage.Store),
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

// Yeni yazı formu
func (homepage Homepage) NewPostForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, ok := helpers.GetCurrentUser(r, homepage.Store)
	if !ok {
		http.Redirect(w, r, "/login?return_url=/profile/new-post", http.StatusSeeOther)
		return
	}
	view, err := template.ParseFiles(helpers.Include("homepage/newpost")...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"CurrentUser": user,
		"Categories":  models.Category{}.GetAll(),
		"Alert":       helpers.GetAlert(w, r, homepage.Store),
	}
	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Yeni yazı kaydet
func (homepage Homepage) NewPostSubmit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, ok := helpers.GetCurrentUser(r, homepage.Store)
	if !ok {
		http.Redirect(w, r, "/login?return_url=/profile/new-post", http.StatusSeeOther)
		return
	}

	title := r.FormValue("blog-title")
	slugStr := slug.Make(title)
	description := r.FormValue("blog-desc")
	categoryID, _ := strconv.Atoi(r.FormValue("blog-category"))
	content := r.FormValue("blog-content")

	// Upload (opsiyonel)
	var pictureURL string
	if err := r.ParseMultipartForm(10 << 20); err == nil {
		if file, header, err := r.FormFile("blog-picture"); err == nil && header != nil && header.Filename != "" {
			if _, statErr := os.Stat("uploads"); os.IsNotExist(statErr) {
				_ = os.MkdirAll("uploads", 0755)
			}
			f, openErr := os.OpenFile("uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if openErr == nil {
				_, _ = io.Copy(f, file)
				pictureURL = "uploads/" + header.Filename
			}
		}
	}
	if pictureURL == "" {
		pictureURL = "uploads/mainpage.jpg"
	}

	models.Post{
		Title:       title,
		Slug:        slugStr,
		Description: description,
		CategoryID:  categoryID,
		Content:     content,
		Picture_url: pictureURL,
		UserID:      user.ID,
	}.Add()

	_ = helpers.SetAlert(w, r, "Yazınız eklendi.", homepage.Store)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// Kullanıcı kendi postunu silebilsin (yorumlar dahil)
func (homepage Homepage) DeleteOwnPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, ok := helpers.GetCurrentUser(r, homepage.Store)
	if !ok {
		http.Redirect(w, r, "/login?return_url=/profile", http.StatusSeeOther)
		return
	}

	idStr := params.ByName("id")
	pid, _ := strconv.Atoi(idStr)
	post := models.Post{}.Get(pid)
	if post.ID == 0 || post.UserID != user.ID {
		_ = helpers.SetAlert(w, r, "Bu içeriği silme yetkiniz yok.", homepage.Store)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	// Önce yorumları sil, sonra postu sil
	_ = models.Comment{}.DeleteByPostID(post.ID)
	post.Delete()
	_ = helpers.SetAlert(w, r, "Yazınız ve yorumları silindi.", homepage.Store)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// Kategoriye göre listeleme
func (homepage Homepage) CategoryList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Kategori sayıları
	counts := models.Post{}.CountsByCategory()

	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string { return models.Category{}.Get(categoryID).Title },
		"getDate":     func(t time.Time) string { return fmt.Sprintf("%02d.%02d.%d", t.Day(), int(t.Month()), t.Year()) },
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
		"getAuthor": func(userID uint) adminmodels.User {
			if userID == 0 {
				return adminmodels.User{}
			}
			return adminmodels.User{}.Get("id = ?", userID)
		},
		"isAdminPost": func(userID uint) bool { return userID == 0 },
		"firstLetter": func(s string) string {
			if s == "" {
				return "?"
			}
			r := []rune(s)
			return strings.ToUpper(string(r[0]))
		},
		"catCount": func(catID uint) int { return counts[int(catID)] },
		"totalCount": func() int {
			sum := 0
			for _, v := range counts {
				sum += v
			}
			return sum
		},
	}).ParseFiles(helpers.Include("homepage/list")...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slugStr := params.ByName("slug")
	cat := models.Category{}.Get("slug = ?", slugStr)
	if cat.ID == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{}
	// SADECE onaylı veya admin yazıları göster (kategoriye göre)
	data["Posts"] = models.Post{}.GetAll("(is_approved = ? OR user_id = 0) AND category_id = ?", true, cat.ID)
	data["Categories"] = models.Category{}.GetAll()
	data["ActiveCategorySlug"] = slugStr
	data["Alert"] = helpers.GetAlert(w, r, homepage.Store)
	if user, ok := helpers.GetCurrentUser(r, homepage.Store); ok {
		data["CurrentUser"] = user
	}
	data["ReturnURL"] = r.URL.RequestURI()

	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Kullanıcı kendi postunu düzenleyebilsin (form)
func (homepage Homepage) EditOwnPostForm(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, ok := helpers.GetCurrentUser(r, homepage.Store)
	if !ok {
		http.Redirect(w, r, "/login?return_url=/profile", http.StatusSeeOther)
		return
	}
	idStr := params.ByName("id")
	pid, _ := strconv.Atoi(idStr)
	post := models.Post{}.Get(pid)
	if post.ID == 0 || post.UserID != user.ID {
		_ = helpers.SetAlert(w, r, "Bu içeriği düzenleme yetkiniz yok.", homepage.Store)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
	view, err := template.New("editpost").Funcs(template.FuncMap{
		"safeHTML": func(html string) template.HTML { return template.HTML(html) },
	}).ParseFiles(helpers.Include("homepage/editpost")...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"CurrentUser": user,
		"Post":        post,
		"Categories":  models.Category{}.GetAll(),
		"Alert":       helpers.GetAlert(w, r, homepage.Store),
	}
	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Kullanıcı kendi postunu düzenleyebilsin (submit)
func (homepage Homepage) EditOwnPostSubmit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, ok := helpers.GetCurrentUser(r, homepage.Store)
	if !ok {
		http.Redirect(w, r, "/login?return_url=/profile", http.StatusSeeOther)
		return
	}
	idStr := params.ByName("id")
	pid, _ := strconv.Atoi(idStr)
	post := models.Post{}.Get(pid)
	if post.ID == 0 || post.UserID != user.ID {
		_ = helpers.SetAlert(w, r, "Bu içeriği düzenleme yetkiniz yok.", homepage.Store)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	title := r.FormValue("blog-title")
	description := r.FormValue("blog-desc")
	categoryID, _ := strconv.Atoi(r.FormValue("blog-category"))
	content := r.FormValue("blog-content")

	// Opsiyonel kapak güncelleme
	if err := r.ParseMultipartForm(10 << 20); err == nil {
		if file, header, err := r.FormFile("blog-picture"); err == nil && header != nil && header.Filename != "" {
			if _, statErr := os.Stat("uploads"); os.IsNotExist(statErr) {
				_ = os.MkdirAll("uploads", 0755)
			}
			f, openErr := os.OpenFile("uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if openErr == nil {
				_, _ = io.Copy(f, file)
				post.Picture_url = "uploads/" + header.Filename
			}
		}
	}

	// Güncelle (slug ve isApproved değişmeden kalsın)
	post.Updates(models.Post{
		Title:       title,
		Description: description,
		CategoryID:  categoryID,
		Content:     content,
		Picture_url: post.Picture_url,
	})

	_ = helpers.SetAlert(w, r, "Yazınız güncellendi.", homepage.Store)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
