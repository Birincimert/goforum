package controllers

import (
	"goblog/site/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandleContactForm, iletişim formu gönderimini işler.
func HandleContactForm(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Formdan gelen verileri alalım
	name := r.FormValue("ct-name")
	email := r.FormValue("ct-mail")
	phone := r.FormValue("ct-phone")
	topic := r.FormValue("ct-category")
	message := r.FormValue("ct-desc")

	// Contact modelini kullanarak veritabanına kaydet
	newContact := models.Contact{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Topic:   topic,
		Message: message,
	}
	newContact.Add()

	log.Println("Form başarıyla veritabanına kaydedildi.")
	http.Redirect(w, r, "/contact?success=true", http.StatusFound)
}
