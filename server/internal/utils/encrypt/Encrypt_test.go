package encrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptHash(t *testing.T) {
	t.Run("正常加密验证", func(t *testing.T) {
		password := "123456"
		hashPassword, err := BcryptHash(password)
		assert.Nil(t, err, "加密过程不应该有错误")
		assert.NotEmpty(t, hashPassword, "加密后的密码不应为空")

		// 验证加密结果
		result := BcryptCheck(password, hashPassword)
		assert.True(t, result, "密码验证应该通过")
	})

	t.Run("空字符串测试", func(t *testing.T) {
		hash, err := BcryptHash("")
		assert.Nil(t, err)
		assert.Empty(t, hash)
	})
}

func TestBcryptCheck(t *testing.T) {
	t.Run("错误密码验证", func(t *testing.T) {
		password := "123456"
		hashPassword, _ := BcryptHash(password)
		result := BcryptCheck("wrongpassword", hashPassword)
		assert.False(t, result, "错误密码不应通过验证")
	})
}

func TestMD5(t *testing.T) {
	t.Run("基本MD5测试", func(t *testing.T) {
		assert.Equal(t, "e10adc3949ba59abbe56e057f20f883e", MD5("123456"), "MD5加密结果不正确")
	})

	t.Run("空字符串MD5测试", func(t *testing.T) {
		assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", MD5(""), "空字符串MD5加密结果不正确")
	})
}
