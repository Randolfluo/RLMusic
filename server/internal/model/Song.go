package model

import (
	"gorm.io/gorm"
)

// Song 歌曲模型
type Song struct {
	ID int `gorm:"primaryKey;auto_increment"`
	// 权限控制
	OwnerID    *int   `gorm:"index" json:"owner_id"` // 上传者/所有者ID
	Permission string // public/private

	// 基础信息
	Title string `gorm:"type:varchar(255);index" json:"title"` // 歌名

	// 关联信息
	ArtistID *int   `gorm:"index" json:"artist_id"`                                                      // 关联艺术家ID
	Artist   Artist `gorm:"foreignKey:ArtistID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"` // 关联艺术家模型

	AlbumID *int  `gorm:"index" json:"album_id"`                                                      // 关联专辑ID
	Album   Album `gorm:"foreignKey:AlbumID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"` // 关联专辑模型

	TrackNum int    `json:"track_num"`                    // 轨道号
	DiscNum  int    `json:"disc_num"`                     // 碟号
	Year     string `gorm:"type:varchar(20)" json:"year"` // 年份

	// 文件信息
	FilePath string `gorm:"type:varchar(500);not null" json:"file_path"`   // 存储路径
	FileName string `gorm:"type:varchar(255)" json:"file_name"`            // 文件名
	FileSize int64  `json:"file_size"`                                     // 文件大小(字节)
	Format   string `gorm:"type:varchar(20);default:'flac'" json:"format"` // 格式，默认 flac

	// 音频参数
	Duration   float64 `json:"duration"`    // 时长(秒)
	SampleRate int     `json:"sample_rate"` // 采样率(Hz) ex: 44100
	BitDepth   int     `json:"bit_depth"`   // 位深(bits) ex: 16, 24
	Channels   int     `json:"channels"`    // 声道数
	BitRate    int     `json:"bit_rate"`    // 比特率(kbps)

	IsDelete bool `gorm:"default:false" json:"-"`
}

// FindSongByPath 根据路径查找歌曲
func FindSongByPath(db *gorm.DB, path string) (*Song, error) {
	var song Song
	// 使用 Find 并限制数量，避免 First 查不到时打印错误日志
	if err := db.Where("file_path = ?", path).Limit(1).Find(&song).Error; err != nil {
		return nil, err
	}
	if song.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &song, nil
}

// SaveSong 保存或更新歌曲
func SaveSong(db *gorm.DB, song *Song) (bool, error) {
	if song.ID > 0 {
		err := db.Save(song).Error
		return false, err // updated
	}
	err := db.Create(song).Error
	return true, err // created
}

// GetSongByID 根据ID获取歌曲
func GetSongByID(db *gorm.DB, id string) (*Song, error) {
	var song Song
	err := db.First(&song, id).Error
	return &song, err
}

// GetSongsList 获取歌曲列表
func GetSongsList(db *gorm.DB, userID int, page, pageSize int) ([]Song, int64, error) {
	var songs []Song
	var total int64

	// 查询条件：用户自己的 + 公开的所有
	query := db.Model(&Song{}).
		Preload("Artist").
		Preload("Album").
		Where("owner_id = ?", userID).
		Or("permission = ?", "public")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id desc").Offset(offset).Limit(pageSize).Find(&songs).Error; err != nil {
		return nil, 0, err
	}

	return songs, total, nil
}
