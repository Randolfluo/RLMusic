package model

import "gorm.io/gorm"

// SystemSetting 系统设置模型，用于持久化存储Key-Value配置
type SystemSetting struct {
	Key   string `gorm:"primaryKey;type:varchar(100)"`
	Value string `gorm:"type:text"`
}

// GetSystemSetting 获取系统设置
func GetSystemSetting(db *gorm.DB, key string) (string, error) {
	var setting SystemSetting
	// 使用 Find 替代 First 以避免记录 ErrRecordNotFound 日志
	err := db.Where("key = ?", key).Limit(1).Find(&setting).Error
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

// SetSystemSetting 设置系统设置(存在更新，不存在创建)
func SetSystemSetting(db *gorm.DB, key string, value string) error {
	setting := SystemSetting{
		Key:   key,
		Value: value,
	}
	return db.Save(&setting).Error
}

// GetAllSystemSettings 获取所有的系统设置
func GetAllSystemSettings(db *gorm.DB) ([]SystemSetting, error) {
	var settings []SystemSetting
	err := db.Find(&settings).Error
	return settings, err
}
