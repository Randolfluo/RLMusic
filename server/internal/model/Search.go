package model

import (
	"gorm.io/gorm"
)

// SearchSongs 搜索歌曲
func SearchSongs(db *gorm.DB, keyword string, page int, limit int) ([]Song, int64, error) {
	var songs []Song
	var total int64
	offset := (page - 1) * limit

	query := db.Model(&Song{}).Where("title LIKE ? OR artist_name LIKE ? OR album_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 只选择需要的字段以保持简明
	// 注意，Preload会加载完整的关联对象，但我们可以控制主表选择的字段
	// 这里选择基础信息和统计信息
	if err := query.Select("id", "title", "artist_name", "album_name", "cover_id", "duration", "play_count").
		Preload("Cover"). // 封面需要Preload因为是关联表
		Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		return nil, 0, err
	}

	return songs, total, nil
}

// SearchArtists 搜索艺术家
func SearchArtists(db *gorm.DB, keyword string, page int, limit int) ([]Artist, int64, error) {
	var artists []Artist
	var total int64
	offset := (page - 1) * limit

	query := db.Model(&Artist{}).Where("name LIKE ?", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&artists).Error; err != nil {
		return nil, 0, err
	}

	return artists, total, nil
}

// SearchAlbums 搜索专辑
func SearchAlbums(db *gorm.DB, keyword string, page int, limit int) ([]Album, int64, error) {
	var albums []Album
	var total int64
	offset := (page - 1) * limit

	query := db.Model(&Album{}).Where("title LIKE ?", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Artist").Limit(limit).Offset(offset).Find(&albums).Error; err != nil {
		return nil, 0, err
	}

	return albums, total, nil
}

// SearchPlaylists 搜索歌单
func SearchPlaylists(db *gorm.DB, keyword string, page int, limit int) ([]Playlist, int64, error) {
	var playlists []Playlist
	var total int64
	offset := (page - 1) * limit

	// 只搜索公开的歌单
	query := db.Model(&Playlist{}).
		Where("(title LIKE ? OR description LIKE ?) AND is_public = ?", "%"+keyword+"%", "%"+keyword+"%", true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 不需要Preload Songs，保持简明
	if err := query.Select("id", "title", "description", "cover_url", "play_count", "owner_id").
		Limit(limit).Offset(offset).Find(&playlists).Error; err != nil {
		return nil, 0, err
	}

	return playlists, total, nil
}
