package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenExpired     = errors.New("token 已过期, 请重新登录")
	ErrTokenNotValidYet = errors.New("token 无效, 请重新登录")
	ErrTokenMalformed   = errors.New("token 不正确, 请重新登录")
	ErrTokenInvalid     = errors.New("这不是一个 token, 请重新登录")
)

// MyClaims 定义JWT的声明结构
type MyClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenToken 生成JWT token
func GenToken(secret, issuer string, expireHour, userId int) (string, error) {
	// 创建声明
	claims := MyClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expireHour))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 使用HS256算法创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(secret))
}

// ParseToken 解析JWT token
func ParseToken(secret, tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		// 根据不同的错误类型返回相应的错误信息
		if vErr, ok := err.(*jwt.ValidationError); ok {
			switch {
			case vErr.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, ErrTokenMalformed
			case vErr.Errors&jwt.ValidationErrorExpired != 0:
				return nil, ErrTokenExpired
			case vErr.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, ErrTokenNotValidYet
			default:
				return nil, ErrTokenInvalid
			}
		}
		return nil, ErrTokenInvalid
	}

	// 验证token
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
