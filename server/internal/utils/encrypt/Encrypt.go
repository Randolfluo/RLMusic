package encrypt

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 使用 bcrypt 对字符串进行加密生成一个哈希值
// str: 需要加密的字符串
// 返回加密后的哈希值和可能的错误
func BcryptHash(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(bytes), err
}

// BcryptCheck 使用 bcrypt 对比明文字符串和哈希值
// plain: 明文字符串
// hash: 哈希值
// 返回是否匹配
func BcryptCheck(plain, hash string) bool {
	if plain == "" || hash == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}

// MD5 计算字符串的 MD5 值
// str: 需要计算的字符串
// b: 可选的额外字节
// 返回 MD5 哈希值的十六进制字符串
func MD5(str string, b ...byte) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(b))
}
