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

type SimpleSongResponse struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	ArtistName string  `json:"artist_name"` // from Artist.Name
	AlbumTitle string  `json:"album_title"` // from Album.Title
	Duration   float64 `json:"duration"`
}

type PlaylistResponse struct {
	ID          int                  `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	IsPublic    bool                 `json:"is_public"`
	OwnerID     int                  `json:"owner_id"`
	Songs       []SimpleSongResponse `json:"songs"`
}

// GetPublicPlaylists 获取所有公共歌单
func GetPublicPlaylists(db *gorm.DB) ([]PlaylistResponse, error) {
	var playlists []Playlist
	// 只查询 public
	if err := db.Preload("Songs").Preload("Songs.Artist").Preload("Songs.Album").
		Where("is_public = ?", true).
		Find(&playlists).Error; err != nil {
		return nil, err
	}
	return convertToResponse(playlists), nil
}

// GetPrivatePlaylists 获取用户私有歌单
func GetPrivatePlaylists(db *gorm.DB, userID int) ([]PlaylistResponse, error) {
	var playlists []Playlist
	// 只查询 private 且 owner_id = userID
	if err := db.Preload("Songs").Preload("Songs.Artist").Preload("Songs.Album").
		Where("is_public = ? AND owner_id = ?", false, userID).
		Find(&playlists).Error; err != nil {
		return nil, err
	}
	return convertToResponse(playlists), nil
}

func convertToResponse(playlists []Playlist) []PlaylistResponse {
	var resp []PlaylistResponse
	for _, p := range playlists {
		var songs []SimpleSongResponse
		for _, s := range p.Songs {
			songs = append(songs, SimpleSongResponse{
				ID:         s.ID,
				Title:      s.Title,
				ArtistName: s.ArtistName,
				AlbumTitle: s.AlbumName,
				Duration:   s.Duration,
			})
		}
		resp = append(resp, PlaylistResponse{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			IsPublic:    p.IsPublic,
			OwnerID:     p.OwnerID,
			Songs:       songs,
		})
	}
	return resp
}

// GetUserPlaylists 获取用户可见的歌单 (Legacy)
func GetUserPlaylists(db *gorm.DB, userID int) ([]PlaylistResponse, error) {
	var playlists []Playlist
	if err := db.Preload("Songs").Preload("Songs.Artist").Preload("Songs.Album").
		Where("owner_id = ?", userID).
		Or("is_public = ?", true).
		Find(&playlists).Error; err != nil {
		return nil, err
	}
	return convertToResponse(playlists), nil
}
