package models

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// About, "hakkında" sayfasının içeriğini temsil eder.
type About struct {
	gorm.Model
	Title   string `gorm:"size:255;not null"`
	Content string `gorm:"type:text"`
}

// About modelinin veritabanı tablosunu oluşturur.
func (about About) Migrate() {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.AutoMigrate(&about)
}

// Hakkında sayfasındaki içeriği veritabanına kaydeder veya günceller.
func (about About) Save() {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	var first About
	// İlk kaydı bul, yoksa oluştur.
	db.FirstOrCreate(&first, About{})

	// İçeriği güncelle ve kaydet.
	first.Title = about.Title
	first.Content = about.Content
	db.Save(&first)
}

// Veritabanından "hakkında" içeriğini çeker.
func (about About) Get() About {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return About{}
	}

	var retrievedAbout About
	db.First(&retrievedAbout)
	return retrievedAbout
}
