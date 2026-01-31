package model

import "gorm.io/gorm"

// UpdateUserCreateFileStatus 更新用户文件初始化状态
func UpdateUserCreateFileStatus(db *gorm.DB, username string) error {
	return db.Model(&User{}).Where("username = ? AND is_delete = ?", username, false).Update("is_create_file", true).Error
}
