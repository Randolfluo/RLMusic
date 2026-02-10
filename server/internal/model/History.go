package model

import (
	"time"

	"gorm.io/gorm"
)

type History struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    int       `gorm:"index" json:"user_id"`
	SongID    int       `gorm:"index" json:"song_id"`
	CreatedAt time.Time `json:"created_at"`
	Song      Song      `json:"song"`
}

// AddHistory 添加或更新播放历史
func AddHistory(db *gorm.DB, userID int, songID int) error {
	var history History
	// 使用 Find 替代 First，避免记录不存在时打印错误日志
	result := db.Where("user_id = ? AND song_id = ?", userID, songID).Limit(1).Find(&history)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		// Update existing record
		history.CreatedAt = time.Now()
		return db.Save(&history).Error
	}

	// Create new record
	history = History{
		UserID:    userID,
		SongID:    songID,
		CreatedAt: time.Now(),
	}
	return db.Create(&history).Error
}

// GetUserHistory 获取用户播放历史
func GetUserHistory(db *gorm.DB, userID int, page int, limit int) ([]History, int64, error) {
	var histories []History
	var total int64

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	err := db.Model(&History{}).Preload("Song").Preload("Song.Artist").Preload("Song.Album").Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Preload("Song").
		Preload("Song.Artist").
		Preload("Song.Album").
		Preload("Song.Cover").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&histories).Error

	return histories, total, err
}

// ClearUserHistory 清空用户播放历史
func ClearUserHistory(db *gorm.DB, userID int) error {
	return db.Where("user_id = ?", userID).Delete(&History{}).Error
}
