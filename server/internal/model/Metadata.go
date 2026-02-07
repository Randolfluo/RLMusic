package model

import (
	"time"

	"gorm.io/gorm"
)

// Artist 艺术家模型
type Artist struct {
	ID int `gorm:"primaryKey;auto_increment" json:"id"`

	Name        string `gorm:"type:varchar(255);not null;index" json:"name"` // 艺术家名称
	Description string `gorm:"type:text" json:"description"`                 // 简介
	CoverSongID *int   `json:"cover_song_id"`                                // 封面对应的歌曲ID
	Cover       string `gorm:"-" json:"cover"`                               // 封面URL (计算字段)
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

// Cover 封面图片模型
type Cover struct {
	ID int `gorm:"primaryKey;auto_increment" json:"id"`

	// 核心属性
	Hash string `gorm:"type:varchar(64);uniqueIndex" json:"hash"` // 图片内容的哈希值(MD5/SHA256)，用于去重
	Path string `gorm:"type:varchar(500)" json:"path"`            // 图片存储路径(相对路径)

	// 图片元数据
	Format string `gorm:"type:varchar(10)" json:"format"` // 格式: jpg, png
	Size   int64  `json:"size"`                           // 文件大小
	Width  int    `json:"width"`                          // 宽
	Height int    `json:"height"`                         // 高
}

// FindOrCreateCover 根据Hash查找或创建封面
func FindOrCreateCover(db *gorm.DB, hash string, path string, format string, size int64, width, height int) (*Cover, error) {
	var cover Cover
	// 尝试根据 Hash 查找，使用 Find 避免 First 找不到时打印日志
	if err := db.Where("hash = ?", hash).Limit(1).Find(&cover).Error; err != nil {
		return nil, err
	}

	if cover.ID == 0 {
		// 创建新封面
		cover = Cover{
			Hash:   hash,
			Path:   path,
			Format: format,
			Size:   size,
			Width:  width,
			Height: height,
		}
		if err := db.Create(&cover).Error; err != nil {
			return nil, err
		}
	}
	return &cover, nil
}
