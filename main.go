package main

import (
	admin_models "goforum/admin/models"
	"goforum/config"
	site_models "goforum/site/models"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	// 1. .env dosyasını yükle
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 2. SESSION_KEY'i environment'dan al
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("SESSION_KEY not found in .env file")
	}

	// 3. Konfigürasyonun oluşturma
	store := sessions.NewCookieStore([]byte(sessionKey))

	admin_models.Post{}.Migrate()
	admin_models.User{}.Migrate()
	admin_models.Category{}.Migrate()
	site_models.About{}.Migrate()
	site_models.Contact{}.Migrate()
	site_models.Comment{}.Migrate()
	site_models.CommentVote{}.Migrate()
	site_models.AuthorApplication{}.Migrate() // Author application table (pending applications)

	http.ListenAndServe(":8080", config.Routes(store))
}
