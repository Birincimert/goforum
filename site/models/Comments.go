package models

import (
	"errors"
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID          uint      `gorm:"not null"`
	Name            string    `gorm:"type:varchar(255);not null"`
	Content         string    `gorm:"type:text;not null"`
	PostID          uint      `gorm:"not null"` // Yorumun ait olduğu blog yazısının ID'si
	ParentCommentID *uint     // Hangi yoruma yanıt verildiğini tutar. NULL olabilir.
	Replies         []Comment `gorm:"foreignKey:ParentCommentID"` // Yanıtları tutar
	LikesCount      int64     `gorm:"-"`                          // SSR için beğeni sayısı
}

func (comment Comment) Migrate() error {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		log.Println("Veritabanı bağlantısı kurulamadı:", err)
		return errors.New("database connection failed")
	}

	// AutoMigrate'in kendi hata dönüşü vardır
	if err := db.AutoMigrate(&comment); err != nil {
		log.Println("Yorum tablosu oluşturulurken hata:", err)
		return errors.New("comment table migration failed")
	}
	return nil
}

// Add, yeni bir yorumu veritabanına ekler ve hata döndürür.
func (comment Comment) Add() error {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		log.Println("Veritabanı bağlantısı kurulamadı:", err)
		return errors.New("database connection failed")
	}

	result := db.Create(&comment)
	// Create işleminde hata olup olmadığını kontrol et
	if result.Error != nil {
		log.Println("Yorum eklenirken hata oluştu:", result.Error)
		return result.Error
	}
	return nil
}

// getRepliesRecursive: Bir yorumun altındaki tüm yanıtları rekürsif olarak doldurur
func (comment Comment) getRepliesRecursive(db *gorm.DB, postID uint, parentID uint) ([]Comment, error) {
	var replies []Comment
	if err := db.Model(&Comment{}).
		Select("comments.*, (SELECT COUNT(1) FROM comment_votes v WHERE v.comment_id = comments.id AND v.value = 1) AS likes_count").
		Where("post_id = ? AND parent_comment_id = ?", postID, parentID).
		Order("created_at ASC").
		Find(&replies).Error; err != nil {
		return nil, err
	}
	for i := range replies {
		nested, err := comment.getRepliesRecursive(db, postID, replies[i].ID)
		if err == nil {
			replies[i].Replies = nested
		}
	}
	return replies, nil
}

func (comment Comment) GetAllByPostID(postID uint) ([]Comment, error) {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Ana yorumları al (likes_count ile)
	var comments []Comment
	if err := db.Model(&Comment{}).
		Select("comments.*, (SELECT COUNT(1) FROM comment_votes v WHERE v.comment_id = comments.id AND v.value = 1) AS likes_count").
		Where("post_id = ? AND parent_comment_id IS NULL", postID).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, err
	}

	// Her bir ana yorum için doğrudan yanıtları al
	for i := range comments {
		replies, err := comment.getRepliesRecursive(db, postID, comments[i].ID)
		if err == nil {
			comments[i].Replies = replies
		}
	}

	return comments, nil
}

// GetAllComments tüm yorumları veritabanından çeker
func (comment Comment) GetAllComments() ([]Comment, error) {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var comments []Comment
	result := db.Order("created_at desc").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}

	return comments, nil
}

// Yorumu sil (yanıtlarıyla birlikte)
func (comment Comment) Delete(id uint) error {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		return err
	}

	// Rekürsif soft delete
	var deleteRecursive func(parentID uint) error
	deleteRecursive = func(parentID uint) error {
		// Önce çocukları bul
		var children []Comment
		if err := db.Where("parent_comment_id = ?", parentID).Find(&children).Error; err != nil {
			return err
		}
		// Her çocuk için altlarını sil
		for _, ch := range children {
			if err := deleteRecursive(ch.ID); err != nil {
				return err
			}
		}
		// En sonda kendisini sil (soft delete)
		if err := db.Delete(&Comment{}, parentID).Error; err != nil {
			return err
		}
		return nil
	}

	return deleteRecursive(id)
}
