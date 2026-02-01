package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

var AesKey = []byte("12345678901234567890123456789012") // 32 bytes for AES-256
var AesIv = []byte("1234567890123456")                  // 16 bytes IV

// AesDecrypt AES解密 (CBC模式, PKCS7填充)
func AesDecrypt(ciphertextBase64 string) (string, error) {
	// Base64 Decode
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(AesKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// CBC mode
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, AesIv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// PKCS7 Unpadding
	plaintext, err = PKCS7Unpadding(plaintext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func PKCS7Unpadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("pkcs7 unpadding error: input empty")
	}
	unpadding := int(origData[length-1])
	if length < unpadding {
		return nil, errors.New("pkcs7 unpadding error: invalid padding")
	}
	return origData[:(length - unpadding)], nil
}
