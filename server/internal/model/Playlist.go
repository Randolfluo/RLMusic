package model

import (
	"strings"

	"gorm.io/gorm"
)

// Playlist 歌单模型
type Playlist struct {
	ID int `gorm:"primaryKey;auto_increment" json:"id"`

	Title       string `gorm:"type:varchar(255);not null;index" json:"title"` // 歌单标题
	Description string `gorm:"type:text" json:"description"`                  // 歌单描述
	IsPublic    bool   `gorm:"index" json:"is_public"`                        // 是否公开 (GORM Default tag removed to allow false/zero-value insert)
	CoverUrl    string `gorm:"type:varchar(500)" json:"cover_url"`            // 歌单封面

	OwnerID int `gorm:"index" json:"owner_id"` // 创建者ID

	Songs []Song `gorm:"many2many:playlist_songs;" json:"songs"` // 歌单包含的歌曲 (多对多)

	// 统计信息
	PlayCount  int `gorm:"default:0" json:"play_count"`  // 播放次数
	TotalSongs int `gorm:"default:0" json:"total_songs"` // 歌曲总数
}

// FindOrCreatePlaylist 查找或创建歌单
func FindOrCreatePlaylist(db *gorm.DB, userID int, title string) (*Playlist, error) {
	var playlist Playlist
	// 使用 Find 避免 First 在找不到记录时打印错误日志
	if err := db.Where("title = ? AND owner_id = ?", title, userID).Limit(1).Find(&playlist).Error; err != nil {
		return nil, err
	}

	// 强制为公开
	isPublic := true

	// ID 为 0 说明未找到
	if playlist.ID == 0 {
		// 如果不存在歌单，则创建
		playlist = Playlist{
			Title:       title,
			Description: "", // 移除权限描述
			IsPublic:    isPublic,
			OwnerID:     userID,
		}
		if err := db.Create(&playlist).Error; err != nil {
			return nil, err
		}
	}
	// else {
	// 	// 如果已存在，不再强制覆盖 is_public，尊重用户的修改
	// }
	return &playlist, nil
}

// CreatePlaylist 创建歌单
func CreatePlaylist(db *gorm.DB, userID int, title string, description string, isPublic bool) (*Playlist, error) {
	playlist := Playlist{
		Title:       title,
		Description: description,
		IsPublic:    isPublic,
		OwnerID:     userID,
	}
	if err := db.Create(&playlist).Error; err != nil {
		return nil, err
	}
	return &playlist, nil
}

// AddSongToPlaylist 添加歌曲到歌单
func AddSongToPlaylist(db *gorm.DB, playlist *Playlist, song *Song) error {
	if err := db.Model(playlist).Association("Songs").Append(song); err != nil {
		return err
	}
	// Update count
	return UpdatePlaylistSongCount(db, playlist)
}

// AddSongsToPlaylist 批量添加歌曲到歌单
func AddSongsToPlaylist(db *gorm.DB, playlist *Playlist, songs []Song) error {
	if err := db.Model(playlist).Association("Songs").Append(songs); err != nil {
		return err
	}
	// Update count
	return UpdatePlaylistSongCount(db, playlist)
}

// RemoveSongFromPlaylist 从歌单移除歌曲
func RemoveSongFromPlaylist(db *gorm.DB, playlist *Playlist, song *Song) error {
	if err := db.Model(playlist).Association("Songs").Delete(song); err != nil {
		return err
	}
	// Update count
	return UpdatePlaylistSongCount(db, playlist)
}

// UpdatePlaylistSongCount 更新歌单歌曲数量
func UpdatePlaylistSongCount(db *gorm.DB, playlist *Playlist) error {
	count := db.Model(playlist).Association("Songs").Count()
	return db.Model(playlist).Update("total_songs", count).Error
}

// IsSongInPlaylist 检查歌曲是否在歌单中
func IsSongInPlaylist(db *gorm.DB, playlistID int, songID int) bool {
	var count int64
	// Table name defaults to playlist_songs
	db.Table("playlist_songs").Where("playlist_id = ? AND song_id = ?", playlistID, songID).Count(&count)
	return count > 0
}

type SimpleSongResponse struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	ArtistName string  `json:"artist_name"` // from Artist.Name
	AlbumTitle string  `json:"album_title"` // from Album.Title
	AlbumName  string  `json:"album_name"`  // Alias for AlbumTitle
	Duration   float64 `json:"duration"`
	Year       string  `json:"year"` // Added Year field

	// 补充字段适配前端播放
	ArtistID int    `json:"artist_id"`
	AlbumID  int    `json:"album_id"`
	CoverUrl string `json:"cover_url"`
}

type PlaylistResponse struct {
	ID          int                  `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	IsPublic    bool                 `json:"is_public"`
	OwnerID     int                  `json:"owner_id"`
	CoverUrl    string               `json:"cover_url"` // 新增封面字段返回
	PlayCount   int                  `json:"play_count"`
	Total       int64                `json:"total"`       // 歌曲总数 (Legacy)
	TotalSongs  int                  `json:"total_songs"` // 歌曲总数
	Songs       []SimpleSongResponse `json:"songs"`       // Deprecated: 列表接口不再返回详情
}

// GetPublicPlaylists 获取所有公开歌单(不含歌曲详情)
func GetPublicPlaylists(db *gorm.DB, page int, limit int) ([]PlaylistResponse, int64, error) {
	var playlists []Playlist
	var total int64
	offset := (page - 1) * limit

	query := db.Model(&Playlist{}).Where("is_public = ?", true)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&playlists).Error; err != nil {
		return nil, 0, err
	}
	// 列表不需要总条数
	return convertToResponse(playlists), total, nil
}

// GetUserPrivatePlaylists 获取用户私有歌单
func GetUserPrivatePlaylists(db *gorm.DB, userID int, page int, limit int) ([]PlaylistResponse, int64, error) {
	var playlists []Playlist
	var total int64
	offset := (page - 1) * limit

	query := db.Model(&Playlist{}).Where("owner_id = ? AND is_public = ?", userID, false)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&playlists).Error; err != nil {
		return nil, 0, err
	}
	return convertToResponse(playlists), total, nil
}

// GetSubscribedPlaylists 获取用户收藏的歌单
func GetSubscribedPlaylists(db *gorm.DB, userID int) ([]PlaylistResponse, error) {
	var user User
	user.ID = userID
	var playlists []Playlist

	if err := db.Model(&user).Association("SubscribedPlaylists").Find(&playlists); err != nil {
		return nil, err
	}
	return convertToResponse(playlists), nil
}

// GetPlaylistDetail 获取完整歌单详情(含歌曲)
func GetPlaylistDetail(db *gorm.DB, playlistIDStr string, page int, limit int) (*PlaylistResponse, error) {
	var playlist Playlist
	// 1. 获取歌单基本信息
	if err := db.First(&playlist, playlistIDStr).Error; err != nil {
		return nil, err
	}

	// 增加播放计数
	db.Model(&playlist).UpdateColumn("play_count", gorm.Expr("play_count + ?", 1))
	playlist.PlayCount++

	// 2. 获取该歌单下的歌曲总数
	total := db.Model(&playlist).Association("Songs").Count()

	// 3. 分页查询歌曲
	var songsRaw []Song
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// 关联查询: Song join playlist_songs
	// 确保表名正确: default GORM many2many table name is playlist_songs
	err := db.Joins("JOIN playlist_songs ON playlist_songs.song_id = song.id").
		Where("playlist_songs.playlist_id = ?", playlist.ID).
		Limit(limit).Offset(offset).
		Preload("Artist").
		Preload("Album").
		Preload("Cover").
		Find(&songsRaw).Error

	if err != nil {
		return nil, err
	}

	if playlist.CoverUrl != "" && !strings.HasPrefix(playlist.CoverUrl, "/covers/") && !strings.HasPrefix(playlist.CoverUrl, "http") {
		playlist.CoverUrl = "/covers/" + playlist.CoverUrl
	}

	var songs []SimpleSongResponse
	for _, s := range songsRaw {
		artistId := 0
		albumId := 0
		if s.ArtistID != nil {
			artistId = *s.ArtistID
		}
		if s.AlbumID != nil {
			albumId = *s.AlbumID
		}

		coverUrl := ""
		// 直接从 Preload 的 Cover 对象获取路径
		if s.CoverID != nil && s.Cover.ID != 0 {
			coverUrl = "/covers/" + s.Cover.Path
		}

		songs = append(songs, SimpleSongResponse{
			ID:         s.ID,
			Title:      s.Title,
			ArtistName: s.ArtistName,
			AlbumTitle: s.AlbumName,
			AlbumName:  s.AlbumName, // Fill AlbumName
			Duration:   s.Duration,
			Year:       s.Year, // Fill Year
			ArtistID:   artistId,
			AlbumID:    albumId,
			CoverUrl:   coverUrl,
		})
	}

	return &PlaylistResponse{
		ID:          playlist.ID,
		Title:       playlist.Title,
		Description: playlist.Description,
		IsPublic:    playlist.IsPublic,
		OwnerID:     playlist.OwnerID,
		CoverUrl:    playlist.CoverUrl,
		PlayCount:   playlist.PlayCount,
		Total:       total,
		TotalSongs:  int(total),
		Songs:       songs,
	}, nil
}

// GetPlaylistRandomSongs 获取歌单随机歌曲(上限100首)
func GetPlaylistRandomSongs(db *gorm.DB, playlistIDStr string, limit int) (*PlaylistResponse, error) {
	var playlist Playlist
	// 1. 获取歌单基本信息
	if err := db.First(&playlist, playlistIDStr).Error; err != nil {
		return nil, err
	}

	// 增加播放计数
	db.Model(&playlist).UpdateColumn("play_count", gorm.Expr("play_count + ?", 1))
	playlist.PlayCount++

	// 2. 获取该歌单下的歌曲总数
	total := db.Model(&playlist).Association("Songs").Count()

	// 3. 随机获取歌曲
	var songsRaw []Song
	if limit > 100 {
		limit = 100
	}

	// 关联查询: Song join playlist_songs
	// SQLite 使用 RANDOM()
	err := db.Joins("JOIN playlist_songs ON playlist_songs.song_id = song.id").
		Where("playlist_songs.playlist_id = ?", playlist.ID).
		Order("RANDOM()").
		Limit(limit).
		Preload("Artist").
		Preload("Album").
		Preload("Cover").
		Find(&songsRaw).Error

	if err != nil {
		return nil, err
	}

	if playlist.CoverUrl != "" && !strings.HasPrefix(playlist.CoverUrl, "/covers/") && !strings.HasPrefix(playlist.CoverUrl, "http") {
		playlist.CoverUrl = "/covers/" + playlist.CoverUrl
	}

	var songs []SimpleSongResponse
	for _, s := range songsRaw {
		artistId := 0
		albumId := 0
		if s.ArtistID != nil {
			artistId = *s.ArtistID
		}
		if s.AlbumID != nil {
			albumId = *s.AlbumID
		}

		coverUrl := ""
		if s.CoverID != nil && s.Cover.ID != 0 {
			coverUrl = "/covers/" + s.Cover.Path
		}

		songs = append(songs, SimpleSongResponse{
			ID:         s.ID,
			Title:      s.Title,
			ArtistName: s.ArtistName,
			AlbumTitle: s.AlbumName,
			AlbumName:  s.AlbumName,
			Duration:   s.Duration,
			Year:       s.Year,
			ArtistID:   artistId,
			AlbumID:    albumId,
			CoverUrl:   coverUrl,
		})
	}

	return &PlaylistResponse{
		ID:          playlist.ID,
		Title:       playlist.Title,
		Description: playlist.Description,
		IsPublic:    playlist.IsPublic,
		OwnerID:     playlist.OwnerID,
		CoverUrl:    playlist.CoverUrl,
		PlayCount:   playlist.PlayCount,
		Total:       total,
		TotalSongs:  int(total),
		Songs:       songs,
	}, nil
}

func convertToResponse(playlists []Playlist) []PlaylistResponse {
	var resp []PlaylistResponse
	for _, p := range playlists {
		// 不再返回歌曲列表
		var songs []SimpleSongResponse = nil // Explicitly nil

		coverUrl := p.CoverUrl
		if coverUrl != "" && !strings.HasPrefix(coverUrl, "/covers/") && !strings.HasPrefix(coverUrl, "http") {
			coverUrl = "/covers/" + coverUrl
		}

		resp = append(resp, PlaylistResponse{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			IsPublic:    p.IsPublic,
			OwnerID:     p.OwnerID,
			CoverUrl:    coverUrl,
			PlayCount:   p.PlayCount,
			TotalSongs:  p.TotalSongs, // 赋值
			Total:       int64(p.TotalSongs),
			Songs:       songs,
		})
	}
	return resp
}

// RemoveSongsFromPlaylist 从歌单中移除歌曲
func RemoveSongsFromPlaylist(db *gorm.DB, playlist *Playlist, songs []Song) error {
	return db.Model(playlist).Association("Songs").Delete(songs)
}

// DeletePlaylist 删除歌单
func DeletePlaylist(db *gorm.DB, playlistID int) error {
	// 1. 删除歌单与歌曲的关联 (GORM 会自动处理，但显式处理更安全或使用 Cascade)
	// 如果使用软删除，关联表记录可能保留。如果物理删除，需清理关联。
	// 这里假设直接删除歌单记录。由于 Many2Many 关系，GORM 默认会删除关联表中的记录。
	return db.Delete(&Playlist{}, playlistID).Error
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
