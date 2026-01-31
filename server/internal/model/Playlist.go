package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Playlist 歌单模型
type Playlist struct {
	ID        int        `gorm:"primaryKey;auto_increment" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	Title       string `gorm:"type:varchar(255);not null;index" json:"title"` // 歌单标题
	Description string `gorm:"type:text" json:"description"`                  // 歌单描述
	IsPublic    bool   `gorm:"default:true" json:"is_public"`                 // 是否公开

	OwnerID int  `gorm:"index" json:"owner_id"`                                                         // 创建者ID
	Owner   User `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"owner"` // 创建者

	CoverImageID *int       `gorm:"index" json:"cover_image_id"`                                                               // 歌单封面
	CoverImage   CoverImage `gorm:"foreignKey:CoverImageID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"cover_image"` // 关联封面

	Songs []Song `gorm:"many2many:playlist_songs;" json:"songs"` // 歌单包含的歌曲 (多对多)
}

// FindOrCreatePlaylist 查找或创建歌单
func FindOrCreatePlaylist(db *gorm.DB, userID int, title string, permission string) (*Playlist, error) {
	var playlist Playlist
	if err := db.Where("title = ? AND owner_id = ?", title, userID).First(&playlist).Error; err != nil {
		// 如果不存在歌单，则创建
		isPublic := permission == "public"
		playlist = Playlist{
			Title:       title,
			Description: fmt.Sprintf("Auto generated %s playlist", permission),
			IsPublic:    isPublic,
			OwnerID:     userID,
		}
		if err := db.Create(&playlist).Error; err != nil {
			return nil, err
		}
	}
	return &playlist, nil
}

// AddSongToPlaylist 添加歌曲到歌单
func AddSongToPlaylist(db *gorm.DB, playlist *Playlist, song *Song) error {
	return db.Model(playlist).Association("Songs").Append(song)
}

// GetUserPlaylists 获取用户可见的歌单
func GetUserPlaylists(db *gorm.DB, userID int) ([]Playlist, error) {
	var playlists []Playlist
	if err := db.Preload("Owner").Where("owner_id = ?", userID).
		Or("is_public = ?", true).
		Find(&playlists).Error; err != nil {
		return nil, err
	}
	return playlists, nil
}
