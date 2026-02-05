package model

import (
	"gorm.io/gorm"
)

type SearchResult struct {
	Songs     []Song     `json:"songs"`
	Artists   []Artist   `json:"artists"`
	Albums    []Album    `json:"albums"`
	Playlists []Playlist `json:"playlists"`
}

// SearchSuggest 搜索建议
func SearchSuggest(db *gorm.DB, keyword string) (*SearchResult, error) {
	var result SearchResult
	limit := 5

	// Search Songs
	if err := db.Where("title LIKE ?", "%"+keyword+"%").Limit(limit).Preload("Artist").Preload("Album").Find(&result.Songs).Error; err != nil {
		return nil, err
	}

	// Search Artists
	if err := db.Where("name LIKE ?", "%"+keyword+"%").Limit(limit).Find(&result.Artists).Error; err != nil {
		return nil, err
	}

	// Search Albums
	if err := db.Where("title LIKE ?", "%"+keyword+"%").Limit(limit).Preload("Artist").Find(&result.Albums).Error; err != nil {
		return nil, err
	}

	// Search Playlists
	if err := db.Where("title LIKE ? AND is_public = ?", "%"+keyword+"%", true).Limit(limit).Find(&result.Playlists).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

// SearchSongs 搜索歌曲
func SearchSongs(db *gorm.DB, keyword string, page int, limit int) ([]Song, int64, error) {
	var songs []Song
	var total int64
	offset := (page - 1) * limit

	query := db.Model(&Song{}).Where("title LIKE ? OR artist_name LIKE ? OR album_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Artist").Preload("Album").Preload("Cover").Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		return nil, 0, err
	}

	return songs, total, nil
}
