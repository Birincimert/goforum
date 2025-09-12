package controllers

import (
	"goforum/admin/helpers"
	"goforum/site/models"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

// Contacts, admin panelinde iletişim formu listeleme işlemini yönetir.
type Contacts struct {
	Store *sessions.CookieStore
}

// Index metodu, tüm iletişim verilerini listeler.
func (c Contacts) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// İletişim mesajlarını veritabanından çek
	allContacts := models.Contact{}.GetAll()
	// Verileri şablon için bir map'e koy
	data := make(map[string]interface{})
	data["Contacts"] = allContacts
	data["Alert"] = helpers.GetAlert(w, r, c.Store)

	// Şablon dosyalarını helper fonksiyonuyla al
	view, err := template.ParseFiles(helpers.Include("contacts")...)
	if err != nil {
		// Hata durumunda loglama yap ve kullanıcıya hata göster
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Şablonu, "index" template adıyla ve verilerle çalıştır
	view.ExecuteTemplate(w, "index", data)
}

// Delete metodu, belirtilen ID'ye sahip iletişim mesajını siler.
func (c Contacts) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	// String olarak gelen ID'yi integer'a dönüştür
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println("ID dönüştürme hatası:", err)
		http.Error(w, "Geçersiz ID", http.StatusBadRequest)
		return
	}

	// Modelden Delete metodunu çağırarak kaydı sil
	models.Contact{}.Delete(idInt)

	// Silme işlemi bittikten sonra kullanıcıyı tekrar listeleme sayfasına yönlendir
	helpers.SetAlert(w, r, "Mesaj başarıyla silindi!", c.Store)
	http.Redirect(w, r, "/admin/contacts", http.StatusSeeOther)
}
