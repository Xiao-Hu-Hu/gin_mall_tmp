package util

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("xiaoxiaohuhexiaoxiaozhou")

type Claims struct {
	ID        uint   `json:"id"`
	UserName  string `json:"username"`
	Authority int    `json:"authority"`
	jwt.RegisteredClaims
}

// GenerateToken 签发 token
func GenerateToken(id uint, userName string, authority int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := Claims{
		ID:        id,
		UserName:  userName,
		Authority: authority,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime), // 令牌过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 令牌签发时间
			NotBefore: jwt.NewNumericDate(time.Now()), // 令牌生效时间
			Issuer:    "xiaoxiaohuhexiaoxiaozhou",     // 签发者标识
		},
	}
	// 创建令牌
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名令牌
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 验证用户 token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

type EmailClaims struct {
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OpeartionType uint   `json:"opeartion_type"`
	Authority     int    `json:"authority"`
	jwt.RegisteredClaims
}

// GenerateEmailToken 签发 email token
func GenerateEmailToken(userID, operation uint, email, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := EmailClaims{
		UserID:        userID,
		Email:         email,
		Password:      password,
		OpeartionType: operation,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime), // 令牌过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 令牌签发时间
			NotBefore: jwt.NewNumericDate(time.Now()), // 令牌生效时间
			Issuer:    "xiaoxiaohuhexiaoxiaozhou",     // 签发者标识
		},
	}
	// 创建令牌
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名令牌
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseEmailToken 验证用户 email token
func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

