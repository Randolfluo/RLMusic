package model

import (
	"fmt"

	"gorm.io/gorm"
)

// Playlist 歌单模型
type Playlist struct {
	ID int `gorm:"primaryKey;auto_increment" json:"id"`

	Title       string `gorm:"type:varchar(255);not null;index" json:"title"` // 歌单标题
	Description string `gorm:"type:text" json:"description"`                  // 歌单描述
	IsPublic    bool   `gorm:"default:true" json:"is_public"`                 // 是否公开

	OwnerID int `gorm:"index" json:"owner_id"` // 创建者ID

	Songs []Song `gorm:"many2many:playlist_songs;" json:"songs"` // 歌单包含的歌曲 (多对多)
}

// FindOrCreatePlaylist 查找或创建歌单
func FindOrCreatePlaylist(db *gorm.DB, userID int, title string, permission string) (*Playlist, error) {
	var playlist Playlist
	// 使用 Find 避免 First 在找不到记录时打印错误日志
	if err := db.Where("title = ? AND owner_id = ?", title, userID).Limit(1).Find(&playlist).Error; err != nil {
		return nil, err
	}

	isPublic := permission == "public"

	// ID 为 0 说明未找到
	if playlist.ID == 0 {
		// 如果不存在歌单，则创建
		playlist = Playlist{
			Title:       title,
			Description: fmt.Sprintf("%s", permission),
			IsPublic:    isPublic,
			OwnerID:     userID,
		}
		if err := db.Create(&playlist).Error; err != nil {
			return nil, err
		}
	} else {
		// 如果已存在，确保 public/private 属性正确 (系统生成的 Public/Private 歌单)
		if playlist.IsPublic != isPublic {
			playlist.IsPublic = isPublic
			// 只更新 is_public 字段
			if err := db.Model(&playlist).Update("is_public", isPublic).Error; err != nil {
				// 忽略错误 or log? handle_song will log if error returned here? No, signature is (*Playlist, error).
				// We can just ignore error here as it's not critical for finding, but good to know.
				return nil, err
			}
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
	if err := db.Where("owner_id = ?", userID).
		Or("is_public = ?", true).
		Find(&playlists).Error; err != nil {
		return nil, err
	}
	return playlists, nil
}
