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
	Password string `gorm:"type:varchar(100)" json:"-"` // bcrypt加密后的密码长度最大为60
	Email    string `gorm:"type:varchar(255)"`

	UserGroup string `gorm:"type:varchar(50);default:'user'" json:"user_group"` // 用户组: admin, user, guest

	Avatar            string `gorm:"type:varchar(255)"` // 头像
	ListeningDuration int64  `gorm:"default:0"`         // 累计听歌时长 (播放产生)
	TotalDuration     int64  `gorm:"default:0"`         // 歌曲总时长 (扫描产生)

	SubscribedPlaylists []Playlist `gorm:"many2many:user_subscribed_playlists;" json:"subscribed_playlists,omitempty"` // 收藏的歌单

	LastLogin *time.Time
	IsDelete  bool
	IPSrc     string `gorm:"type:varchar(50)" json:"ip_src"`
}

// GetUserAuthInfoByName 根据用户名获取用户信息
func GetUserAuthInfoByName(db *gorm.DB, name string) (*User, error) {
	var user User
	// 使用 Limit(1).Find 避免 gorm.First 的 "record not found" 错误日志
	if err := db.Where("username = ? AND is_delete = ?", name, false).Limit(1).Find(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	// slog.Info("GetUserAuthInfoByName", "username", name, "found", true)
	return &user, nil
}

// UpdateUserLoginInfo 更新用户登录信息
func UpdateUserLoginInfo(db *gorm.DB, id int, ip string) error {
	now := time.Now()
	// 更新 LastLogin 和 IPSrc
	return db.Model(&User{}).Where("id = ? AND is_delete = ?", id, false).
		Updates(map[string]interface{}{
			"last_login": now,
			"ip_src":     ip,
		}).Error
}

// UpdateUserAvatar 更新用户头像
func UpdateUserAvatar(db *gorm.DB, id int, avatarPath string) error {
	return db.Model(&User{}).Where("id = ?", id).Update("avatar", avatarPath).Error
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
	// 使用 Limit(1).Find 避免 gorm.First 的 "record not found" 错误日志
	result := db.Where("id = ? AND is_delete = ?", id, false).Limit(1).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
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
	// 使用 Limit(1).Find 避免 gorm.First 的错误日志
	if err := db.Unscoped().Where("username = ?", username).Limit(1).Find(&oldUser).Error; err != nil {
		return err
	}

	if oldUser.ID != 0 {
		timestamp := time.Now().Unix()
		// slog.Info("Renaming soft deleted user", "old_username", username, "id", oldUser.ID)
		// newUsername 已被移除，直接使用 gorm 表达式在数据库侧进行拼接
		return db.Model(&oldUser).Update("username", gorm.Expr("username || '_del_' || ?", timestamp)).Error
	}

	return nil
}

// UpdateListeningDuration 增加用户累计听歌时长
func UpdateListeningDuration(db *gorm.DB, id int, duration int64) error {
	return db.Model(&User{}).Where("id = ?", id).Update("listening_duration", gorm.Expr("listening_duration + ?", duration)).Error
}

// IsPlaylistSubscribed 检查用户是否收藏了指定歌单
func IsPlaylistSubscribed(db *gorm.DB, userID int, playlistID int) (bool, error) {
	var count int64
	// 关联表名为 user_subscribed_playlists (在 User 结构体 tag 中定义)
	err := db.Table("user_subscribed_playlists").
		Where("user_id = ? AND playlist_id = ?", userID, playlistID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
