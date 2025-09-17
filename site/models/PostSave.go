package models

import (
	"errors"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// PostSave: kullanıcının kaydettiği yazıları tutar
// Her (user_id, post_id) çifti benzersizdir.
type PostSave struct {
	gorm.Model
	PostID uint `gorm:"index;not null"`
	UserID uint `gorm:"index;not null"`
}

func (PostSave) Migrate() error {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		return err
	}
	return db.AutoMigrate(&PostSave{})
}

// Toggle: kayıt durumunu değiştirir; sonucunda true ise kaydedildi, false ise kaldırıldı.
func (ps PostSave) Toggle(userID, postID uint) (bool, error) {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		return false, err
	}
	var existing PostSave
	res := db.Where("user_id = ? AND post_id = ?", userID, postID).First(&existing)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// oluştur
			newRow := PostSave{UserID: userID, PostID: postID}
			if err := db.Create(&newRow).Error; err != nil {
				return false, err
			}
			return true, nil
		}
		return false, res.Error
	}
	// vardı -> sil
	if err := db.Delete(&existing).Error; err != nil {
		return false, err
	}
	return false, nil
}

func (ps PostSave) IsSaved(userID, postID uint) (bool, error) {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		return false, err
	}
	var count int64
	if err := db.Model(&PostSave{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetSavedPostsByUser: kullanıcının kaydettiği yazıları getirir (son kaydedilenler önce)
func GetSavedPostsByUser(userID uint) ([]Post, error) {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Önce post id'leri al
	var rows []PostSave
	if err := db.Where("user_id = ?", userID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []Post{}, nil
	}
	ids := make([]uint, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.PostID)
	}
	// Postları sırayı olabildiğince koruyarak getir
	var posts []Post
	if err := db.Where("id IN ?", ids).Find(&posts).Error; err != nil {
		return nil, err
	}
	// Yorumları da doldur
	for i := range posts {
		comments, _ := Comment{}.GetAllByPostID(posts[i].ID)
		posts[i].Comments = comments
	}
	return posts, nil
}
