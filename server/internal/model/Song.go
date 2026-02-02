package model

import (
	"gorm.io/gorm"
)

// Song 歌曲模型
type Song struct {
	ID int `gorm:"primaryKey;auto_increment"`

	// 基础信息
	Title      string `gorm:"type:varchar(255);index" json:"title"`       // 歌名
	ArtistName string `gorm:"type:varchar(255);index" json:"artist_name"` // 冗余反范式化：直接存储艺术家名称
	AlbumName  string `gorm:"type:varchar(255);index" json:"album_name"`  // 冗余反范式化：直接存储专辑名称

	// 关联信息
	ArtistID *int   `gorm:"index" json:"artist_id"`                                                      // 关联艺术家ID
	Artist   Artist `gorm:"foreignKey:ArtistID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"` // 关联艺术家模型

	AlbumID *int  `gorm:"index" json:"album_id"`                                                      // 关联专辑ID
	Album   Album `gorm:"foreignKey:AlbumID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"` // 关联专辑模型

	CoverID *int  `gorm:"index" json:"cover_id"`                                                      // 关联封面ID
	Cover   Cover `gorm:"foreignKey:CoverID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"` // 关联封面模型

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

	// 统计信息
	PlayCount int `gorm:"default:0" json:"play_count"` // 播放次数

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
