package model

import (
	"time"

	"gorm.io/gorm"
)

// Artist 艺术家模型
type Artist struct {
	ID        int        `gorm:"primaryKey;auto_increment" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	Name        string `gorm:"type:varchar(255);not null;index" json:"name"` // 艺术家名称
	Description string `gorm:"type:text" json:"description"`                 // 简介
}

// Album 专辑模型
type Album struct {
	ID        int        `gorm:"primaryKey;auto_increment" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	Title       string     `gorm:"type:varchar(255);not null;index" json:"title"` // 专辑标题
	Description string     `gorm:"type:text" json:"description"`                  // 专辑简介
	ReleaseDate *time.Time `json:"release_date"`                                  // 发行日期

	ArtistID *int   `gorm:"index" json:"artist_id"`                                                           // 关联艺术家ID
	Artist   Artist `gorm:"foreignKey:ArtistID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"artist"` // 关联艺术家

	CoverImageID *int       `gorm:"index" json:"cover_image_id"`                                                               // 关联封面ID
	CoverImage   CoverImage `gorm:"foreignKey:CoverImageID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"cover_image"` // 关联封面
}

// CoverImage 封面图片模型
type CoverImage struct {
	ID        int        `gorm:"primaryKey;auto_increment" json:"id"`
	CreatedAt *time.Time `json:"created_at"`

	Data     []byte `gorm:"type:mediumblob" json:"-"`               // 图片数据 (直接存储)
	MimeType string `gorm:"type:varchar(50)" json:"mime_type"`      // 图片类型 (image/jpeg, image/png)
	Checksum string `gorm:"type:varchar(64);index" json:"checksum"` // 文件哈希，防止重复存储
}

// FindOrCreateArtist 查找或创建艺术家
func FindOrCreateArtist(db *gorm.DB, name string) (*Artist, error) {
	var artist Artist
	if err := db.FirstOrCreate(&artist, Artist{Name: name}).Error; err != nil {
		return nil, err
	}
	return &artist, nil
}

// FindOrCreateCover 查找或创建封面
func FindOrCreateCover(db *gorm.DB, data []byte, mime, checksum string) (*int, error) {
	var cover CoverImage
	if err := db.Where("checksum = ?", checksum).First(&cover).Error; err == nil {
		return &cover.ID, nil
	}

	newCover := CoverImage{
		Data:     data,
		MimeType: mime,
		Checksum: checksum,
	}
	if err := db.Create(&newCover).Error; err != nil {
		return nil, err
	}
	return &newCover.ID, nil
}

// GetCover 获取封面
func GetCover(db *gorm.DB, id int) (*CoverImage, error) {
	var cover CoverImage
	err := db.First(&cover, id).Error
	return &cover, err
}

// FindOrCreateAlbum 查找或创建专辑
func FindOrCreateAlbum(db *gorm.DB, title string, artistID int, coverID *int) (*Album, error) {
	var album Album
	// 尝试查找同名且同艺术家的专辑
	if err := db.Where("title = ? AND artist_id = ?", title, artistID).First(&album).Error; err != nil {
		// 没找到，创建
		album = Album{
			Title:        title,
			ArtistID:     &artistID,
			CoverImageID: coverID,
		}
		if err := db.Create(&album).Error; err != nil {
			return nil, err
		}
	}
	return &album, nil
}
