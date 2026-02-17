package model

import (
	"strconv"

	"gorm.io/gorm"
)

const (
	KeyTotalSongs    = "total_songs"
	KeyTotalAlbums   = "total_albums"
	KeyTotalArtists  = "total_artists"
	KeyTotalDuration = "total_duration"
	KeyTotalVolume   = "total_volume"
)

// SystemInfo 系统信息模型 (Key-Value)
type SystemInfo struct {
	Key   string `gorm:"primaryKey;type:varchar(100)"`
	Value string `gorm:"type:text"`
}

// SystemInfoStruct 用于应用层传输的结构体 (非数据库表)
type SystemInfoStruct struct {
	TotalSongs    int64
	TotalAlbums   int64
	TotalArtists  int64
	TotalDuration int64
	TotalVolume   int64
}

// GetSystemInfoValue 获取单个系统设置
func GetSystemInfoValue(db *gorm.DB, key string) (string, error) {
	var info SystemInfo
	err := db.Where("`key` = ?", key).Limit(1).Find(&info).Error
	if err != nil {
		return "", err
	}
	return info.Value, nil
}

// SetSystemInfoValue 设置单个系统设置
func SetSystemInfoValue(db *gorm.DB, key string, value string) error {
	info := SystemInfo{Key: key, Value: value}
	return db.Save(&info).Error
}

// GetSystemInfoStruct 获取汇总的系统信息结构体
func GetSystemInfoStruct(db *gorm.DB) (*SystemInfoStruct, error) {
	var infos []SystemInfo
	if err := db.Find(&infos).Error; err != nil {
		return nil, err
	}

	res := &SystemInfoStruct{}
	for _, info := range infos {
		switch info.Key {
		case KeyTotalSongs:
			res.TotalSongs, _ = strconv.ParseInt(info.Value, 10, 64)
		case KeyTotalAlbums:
			res.TotalAlbums, _ = strconv.ParseInt(info.Value, 10, 64)
		case KeyTotalArtists:
			res.TotalArtists, _ = strconv.ParseInt(info.Value, 10, 64)
		case KeyTotalDuration:
			res.TotalDuration, _ = strconv.ParseInt(info.Value, 10, 64)
		case KeyTotalVolume:
			res.TotalVolume, _ = strconv.ParseInt(info.Value, 10, 64)
		}
	}
	return res, nil
}

// UpdateSystemInfoStats 更新统计数据
func UpdateSystemInfoStats(db *gorm.DB, songs, albums, artists, duration, volume int64) error {
	if err := SetSystemInfoValue(db, KeyTotalSongs, strconv.FormatInt(songs, 10)); err != nil {
		return err
	}
	if err := SetSystemInfoValue(db, KeyTotalAlbums, strconv.FormatInt(albums, 10)); err != nil {
		return err
	}
	if err := SetSystemInfoValue(db, KeyTotalArtists, strconv.FormatInt(artists, 10)); err != nil {
		return err
	}
	if err := SetSystemInfoValue(db, KeyTotalDuration, strconv.FormatInt(duration, 10)); err != nil {
		return err
	}
	if err := SetSystemInfoValue(db, KeyTotalVolume, strconv.FormatInt(volume, 10)); err != nil {
		return err
	}
	return nil
}
