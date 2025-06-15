package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TokenExpireDuration = time.Hour * 24 * 7 //定义过期时间

var mySecret = []byte("3g") //自定义密码加盐，将这个加盐和原始信息拼接一起加密

type MyClaims struct {
	// 可根据需要自行添加字段
	UserID             int    `json:"user_id"`
	UserName           string `json:"user_name"`
	Status             int    `json:"status"`
	jwt.StandardClaims        // 内嵌标准的声明
}

// 生成JWT
func GenToken(userID int, userName string, Status int) (string, error) {
	claims := MyClaims{
		userID,
		userName,
		Status,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), //过期时间
			Issuer:    "admin",                                    //签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySecret)
}

// 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return mc, nil
}
