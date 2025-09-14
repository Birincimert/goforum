package models

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type AuthorApplication struct {
	gorm.Model
	FirstName string `gorm:"type:nvarchar(50);not null"`
	LastName  string `gorm:"type:nvarchar(50);not null"`
	Email     string `gorm:"type:nvarchar(100);not null;unique"`
	Username  string `gorm:"type:nvarchar(50);not null;unique"`
	Password  string `gorm:"type:nvarchar(255);not null"`
	Bio       string `gorm:"type:nvarchar(500);null"`
	Status    string `gorm:"type:nvarchar(50);not null;default:'pending'"` // pending/approved/rejected
}

func (app AuthorApplication) Migrate() {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := db.AutoMigrate(&app); err != nil {
		fmt.Println(err)
	}
}

func (app AuthorApplication) Add() {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Create(&app)
}

func (app AuthorApplication) Get(where ...interface{}) AuthorApplication {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return app
	}
	db.First(&app, where...)
	return app
}

func (app AuthorApplication) GetAll(where ...interface{}) []AuthorApplication {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var apps []AuthorApplication
	db.Find(&apps, where...)
	return apps
}

func (app AuthorApplication) UpdateStatus(status string) {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Model(&app).Update("Status", status)
}

func (app AuthorApplication) Delete() {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Delete(&app, app.ID)
}
