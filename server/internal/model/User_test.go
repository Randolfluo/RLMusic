package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	return db
}

func TestUserCRUD(t *testing.T) {
	db := setupTestDB()

	t.Run("CreateUser", func(t *testing.T) {
		user, err := CreateUser(db, "testuser", "hashedpassword", "test@example.com")
		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
	})

	t.Run("GetUserByName", func(t *testing.T) {
		user, err := GetUserAuthInfoByName(db, "testuser")
		assert.Nil(t, err)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("UpdateLoginInfo", func(t *testing.T) {
		user, _ := GetUserAuthInfoByName(db, "testuser")
		// oldTime := user.LastLogin // Initialize is null

		err := UpdateUserLoginInfo(db, user.ID)
		assert.Nil(t, err)

		// Re-fetch
		updatedUser, _ := GetUserAuthInfoById(db, user.ID)
		assert.NotNil(t, updatedUser.LastLogin)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		user, _ := GetUserAuthInfoByName(db, "testuser")
		err := DeleteUser(db, user.ID)
		assert.Nil(t, err)

		// Check if deleted
		_, err = GetUserAuthInfoById(db, user.ID)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
