package model

import (
	"gorm.io/gorm"
)

// MakeMigrate 执行数据库迁移
func MakeMigrate(db *gorm.DB) error {
	// 按顺序迁移: 基础表 -> 关联表 -> 主表
	if err := db.AutoMigrate(
		&User{},
		&SystemSetting{},
		&Artist{},
		&Album{},
		&Song{},
		&Playlist{},
	); err != nil {
		return err
	}
	return nil
}
