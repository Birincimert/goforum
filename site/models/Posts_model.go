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
	UserID                                         uint
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
	// Kısmi preload yerine önce postları al, ardından her post için yorumları derin olarak doldur
	db.Order("id DESC").Find(&posts, where...)
	for i := range posts {
		comments, _ := Comment{}.GetAllByPostID(posts[i].ID)
		posts[i].Comments = comments
	}
	return posts
}

// CountsByCategory: kategori_id'ye göre toplam post sayıları
func (post Post) CountsByCategory() map[int]int {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return map[int]int{}
	}
	type Row struct {
		CategoryID int
		Cnt        int64
	}
	var rows []Row
	db.Model(&Post{}).Select("category_id as category_id, COUNT(*) as cnt").Group("category_id").Scan(&rows)
	counts := make(map[int]int, len(rows))
	for _, r := range rows {
		counts[r.CategoryID] = int(r.Cnt)
	}
	return counts
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
