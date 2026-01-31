package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        int `gorm:"primary_key;auto_increment"`
	CreatedAt *time.Time
	UpdatedAt *time.Time

	Username string `gorm:"type:varchar(50);unique"`
	Password string `gorm:"type:varchar(100)"` // bcrypt加密后的密码长度最大为60
	Email    string `gorm:"type:varchar(255)"`

	UserGroup string `gorm:"type:varchar(50);default:'user'" json:"user_group"` // 用户组: admin, user, guest

	LastLogin    *time.Time
	IsDelete     bool
	IsCreateFile bool
	// IPAddr    string    `gorm:"type:varchar(20)" json:"ip_addr"`
	// IPSrc     string    `gorm:"type:varchar(50)" json:"ip_src"`
}

// GetUserAuthInfoByName 根据用户名获取用户信息
func GetUserAuthInfoByName(db *gorm.DB, name string) (*User, error) {
	var user User
	err := db.Where("username = ? AND is_delete = ?", name, false).First(&user).Error
	return &user, err
}

// UpdateUserLoginInfo 更新用户登录信息
func UpdateUserLoginInfo(db *gorm.DB, id int) error {
	now := time.Now()
	return db.Model(&User{}).Where("id = ? AND is_delete = ?", id, false).Update("last_login", now).Error
}

// CreateUser 创建新用户
func CreateUser(db *gorm.DB, username string, password string, email string) (*User, error) {
	now := time.Now()
	user := &User{
		CreatedAt: &now,
		UpdatedAt: &now,
		Username:  username,
		Password:  password,
		Email:     email,
	}

	err := db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserAuthInfoById 根据ID获取用户信息
func GetUserAuthInfoById(db *gorm.DB, id int) (*User, error) {
	var user User
	err := db.Where("id = ? AND is_delete = ?", id, false).First(&user).Error
	return &user, err
}

// DeleteUser 删除用户 (软删除)
// 删除时给用户名加上后缀，释放唯一索引，允许后续重新注册该用户名
func DeleteUser(db *gorm.DB, id int) error {
	timestamp := time.Now().Unix()
	return db.Model(&User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_delete": true,
			"username":  gorm.Expr("username || '_del_' || ?", timestamp),
		}).Error
}

// RenameSoftDeletedUser 重命名已软删除的用户
// 用于处理历史遗留的未释放用户名的软删除记录
func RenameSoftDeletedUser(db *gorm.DB, username string) error {
	var oldUser User
	if err := db.Unscoped().Where("username = ?", username).First(&oldUser).Error; err == nil {
		timestamp := time.Now().Unix()
		// newUsername 已被移除，直接使用 gorm 表达式在数据库侧进行拼接
		return db.Model(&oldUser).Update("username", gorm.Expr("username || '_del_' || ?", timestamp)).Error
	}

	return nil
}
