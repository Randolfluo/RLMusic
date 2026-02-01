package model

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
)

// Artist 艺术家模型
type Artist struct {
	ID int `gorm:"primaryKey;auto_increment" json:"id"`

	Name        string `gorm:"type:varchar(255);not null;index" json:"name"` // 艺术家名称
	Description string `gorm:"type:text" json:"description"`                 // 简介
}

// Album 专辑模型
type Album struct {
	ID int `gorm:"primaryKey;auto_increment" json:"id"`

	Title       string     `gorm:"type:varchar(255);not null;index" json:"title"` // 专辑标题
	Description string     `gorm:"type:text" json:"description"`                  // 专辑简介
	ReleaseDate *time.Time `json:"release_date"`                                  // 发行日期

	ArtistID *int   `gorm:"index" json:"artist_id"`                                                           // 关联艺术家ID
	Artist   Artist `gorm:"foreignKey:ArtistID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"artist"` // 关联艺术家
}

// FindOrCreateArtist 查找或创建艺术家
func FindOrCreateArtist(db *gorm.DB, name string) (*Artist, error) {
	var artist Artist
	if err := db.FirstOrCreate(&artist, Artist{Name: name}).Error; err != nil {
		return nil, err
	}
	return &artist, nil
}

// FindOrCreateAlbum 查找或创建专辑
func FindOrCreateAlbum(db *gorm.DB, title string, artistID *int) (*Album, error) {
	var album Album

	// 构建查询条件
	query := db.Where("title = ?", title)
	if artistID != nil {
		query = query.Where("artist_id = ?", *artistID)
	} else {
		query = query.Where("artist_id IS NULL")
	}

	// 尝试查找同名且同艺术家的专辑，使用 Find 避免日志报错
	if err := query.Limit(1).Find(&album).Error; err != nil {
		return nil, err
	}

	if album.ID == 0 {
		// 没找到，创建
		slog.Info("Creating new album", "title", title)
		album = Album{
			Title:    title,
			ArtistID: artistID,
		}
		if err := db.Create(&album).Error; err != nil {
			return nil, err
		}
	}
	return &album, nil
}
