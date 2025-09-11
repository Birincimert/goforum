package models

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title, Slug, Description, Content, Picture_url string
	CategoryID                                     int
	Comments                                       []Comment
}

func (post Post) Migrate() {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.AutoMigrate(&post)
}

func (post Post) Add() {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Create(&post)
}

//value interface olması her tipten değer gelebilir anlamına gelir.

func (post Post) Get(where ...interface{}) Post {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return post
	}

	// Post'u al
	db.First(&post, where...)

	// Yorumları getir
	comments, _ := Comment{}.GetAllByPostID(post.ID)
	post.Comments = comments

	return post
}

func (post Post) GetBySlug(slug string) Post {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return post
	}

	// Post'u al
	db.First(&post, "slug = ?", slug)

	// Yorumları getir (likes_count dahil)
	comments, _ := Comment{}.GetAllByPostID(post.ID)
	post.Comments = comments

	return post
}

func (post Post) GetAll(where ...interface{}) []Post {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var posts []Post
	db.Preload("Comments", "parent_comment_id IS NULL").
		Preload("Comments.Replies").
		Order("id DESC").
		Find(&posts, where...)
	return posts
}

func (post Post) Update(column string, value interface{}) {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Model(&post).Update(column, value)
}

func (post Post) Updates(data Post) {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Model(&post).Updates(data)
}

func (post Post) Delete() {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Delete(&post, post.ID)
}
