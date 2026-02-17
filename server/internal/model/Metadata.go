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
	Cover       string `gorm:"type:varchar(500)" json:"cover"`               // 封面URL

	Songs []Song `gorm:"many2many:song_artists;" json:"songs,omitempty"`
}

// Album 专辑模型
type Album struct {
	ID int `gorm:"primaryKey;auto_increment" json:"id"`

	Title       string     `gorm:"type:varchar(255);not null;index" json:"title"` // 专辑标题
	Description string     `gorm:"type:text" json:"description"`                  // 专辑简介
	ReleaseDate *time.Time `json:"release_date"`                                  // 发行日期

	Cover string `gorm:"type:varchar(500)" json:"cover"` // 封面URL

	ArtistID *int   `gorm:"index" json:"artist_id"`                                                           // 关联艺术家ID
	Artist   Artist `gorm:"foreignKey:ArtistID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"artist"` // 关联艺术家

	Songs []Song `gorm:"foreignKey:AlbumID" json:"songs,omitempty"` // 专辑包含的歌曲
}

// FindOrCreateArtist 查找或创建艺术家
func FindOrCreateArtist(db *gorm.DB, name string) (*Artist, error) {
	var artist Artist
	if err := db.FirstOrCreate(&artist, Artist{Name: name, Description: ""}).Error; err != nil {
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
			Title:       title,
			ArtistID:    artistID,
			Description: "",
		}
		if err := db.Create(&album).Error; err != nil {
			return nil, err
		}
	}
	return &album, nil
}

// GetArtistRandomSongs 获取艺术家随机歌曲(上限100首)
func GetArtistRandomSongs(db *gorm.DB, artistIDStr string, limit int) ([]SimpleSongResponse, error) {
	if limit > 100 {
		limit = 100
	}

	var songsRaw []Song
	// 关联查询: Song join song_artists
	// SQLite 使用 RANDOM()
	err := db.Joins("JOIN song_artists ON song_artists.song_id = song.id").
		Where("song_artists.artist_id = ?", artistIDStr).
		Order("RANDOM()").
		Limit(limit).
		Preload("Artist").
		Preload("Album").
		Preload("Cover").
		Find(&songsRaw).Error

	if err != nil {
		return nil, err
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
	return songs, nil
}

// GetAlbumRandomSongs 获取专辑随机歌曲(上限100首)
func GetAlbumRandomSongs(db *gorm.DB, albumIDStr string, limit int) ([]SimpleSongResponse, error) {
	if limit > 100 {
		limit = 100
	}

	var songsRaw []Song
	// 专辑ID直接在Song表里
	err := db.Where("album_id = ?", albumIDStr).
		Order("track_num ASC"). // 专辑通常按音轨排序更好，但为了保持"随机"接口语义一致性，或者如果真要随机分析：
		// Order("RANDOM()"). // AI分析通常需要完整专辑感，但如果只想取样...
		// 专辑歌曲通常不多，直接全部取出，前端或调用方自己截取。
		// 但为了兼容大专辑或一致性，我们还是支持limit
		Limit(limit).
		Preload("Artist").
		Preload("Album").
		Preload("Cover").
		Find(&songsRaw).Error

	if err != nil {
		return nil, err
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
	return songs, nil
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
