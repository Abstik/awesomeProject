package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TokenExpireDuration = time.Hour * 24 * 1 //定义过期时间

var mySecret = []byte("xiyou_mobile_lab@2024#jwt!secret_key")

type MyClaims struct {
	// 可根据需要自行添加字段
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Status   int    `json:"status"`
	jwt.RegisteredClaims
}

// 生成JWT
func GenToken(userID int, userName string, Status int) (string, error) {
	claims := MyClaims{
		userID,
		userName,
		Status,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)), //过期时间
			Issuer:    "admin",                                                 //签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySecret)
}

// 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return mc, nil
}
