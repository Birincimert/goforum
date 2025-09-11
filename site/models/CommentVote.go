package models

import (
	"errors"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type CommentVote struct {
	gorm.Model
	CommentID uint `gorm:"index;not null"`
	UserID    uint `gorm:"index;not null"`
	Value     int  `gorm:"not null"` // 1 = like
}

func (CommentVote) Migrate() error {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		return err
	}
	return db.AutoMigrate(&CommentVote{})
}

// SetVote sets or toggles a like for a user on a comment
func (cv CommentVote) SetVote(userID uint, commentID uint, value int) error {
	if value != 1 && value != -1 {
		return errors.New("invalid vote value")
	}
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		return err
	}
	var existing CommentVote
	tx := db.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&existing)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			// create new like
			existing = CommentVote{UserID: userID, CommentID: commentID, Value: 1}
			return db.Create(&existing).Error
		}
		return tx.Error
	}
	// toggle
	return db.Delete(&existing).Error
}

// CountVotes returns like count for a comment
func (cv CommentVote) CountVotes(commentID uint) (int64, error) {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		return 0, err
	}
	var up int64
	if err := db.Model(&CommentVote{}).Where("comment_id = ? AND value = 1", commentID).Count(&up).Error; err != nil {
		return 0, err
	}
	return up, nil
}

// IsLikedBy returns whether a user liked a specific comment
func (cv CommentVote) IsLikedBy(userID uint, commentID uint) (bool, error) {
	db, err := gorm.Open(sqlserver.Open(Dns), &gorm.Config{})
	if err != nil {
		return false, err
	}
	var count int64
	if err := db.Model(&CommentVote{}).Where("user_id = ? AND comment_id = ? AND value = 1", userID, commentID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
