package controllers

import (
	"goforum/admin/helpers"
	site_models "goforum/site/models"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

// About, admin paneldeki "Hakkında" sayfasını yönetir.
type About struct {
	Store *sessions.CookieStore
}

// Index metodu, "hakkında" sayfasını düzenleme formunu gösterir.
func (a About) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Kullanıcı oturumunu kontrol et
	if !helpers.CheckUser(w, r, a.Store) {
		return
	}

	// Veritabanından mevcut içeriği çek. Eğer kayıt yoksa, boş bir struct döner.
	aboutData := site_models.About{}.Get()

	// Şablon dosyalarını yükle
	view, err := template.New("index").ParseFiles(helpers.Include("about")...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["Alert"] = helpers.GetAlert(w, r, a.Store)
	data["About"] = aboutData // About verisini burada map'e ekledik

	// Veriyi şablona gönder
	view.ExecuteTemplate(w, "index", data)
}

// Save metodu, formu işler ve veritabanını günceller.
func (a About) Save(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Kullanıcı oturumunu kontrol et
	if !helpers.CheckUser(w, r, a.Store) {
		return
	}

	r.ParseForm()

	// Formdan gelen verileri modelin içine at
	about := site_models.About{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
	}

	// `Save` metodu ile veriyi veritabanına kaydet veya güncelle
	about.Save()

	// Başarılı kayıttan sonra admin sayfasına yönlendir
	helpers.SetAlert(w, r, "Hakkında sayfası başarıyla güncellendi.", a.Store)
	http.Redirect(w, r, "/admin/about", http.StatusSeeOther)
}
