package models

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Contact, veritabanındaki "contacts" tablosunu temsil eder.
type Contact struct {
	gorm.Model
	Name    string `gorm:"size:255;not null"`
	Email   string `gorm:"size:255;not null"`
	Phone   string `gorm:"size:20"`
	Topic   string `gorm:"size:255"`
	Message string `gorm:"type:text;not null"`
}

// Contact modelinin veritabanı tablosunu oluşturur.
func (contact Contact) Migrate() {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.AutoMigrate(&contact)
}

// Yeni bir iletişim formunu veritabanına kaydeder.
func (contact Contact) Add() {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Create(&contact)
}

// Veritabanından iletişim formunu çekmesi için
func (contact Contact) GetAll() []Contact {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var contacts []Contact
	db.Find(&contacts)
	return contacts
}

// İletişim mesajını veritabanından siler
func (contact Contact) Delete(id int) {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		log.Println("Delete(): Veritabanı bağlantı hatası:", err)
		return
	}

	db.Delete(&contact, id)
}
