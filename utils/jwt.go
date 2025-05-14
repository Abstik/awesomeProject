package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TokenExpireDuration = time.Hour * 24 * 365 //定义过期时间

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
	// 创建一个我们自己的声明
	claims := MyClaims{
		userID,
		userName,
		Status,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), //过期时间
			Issuer:    "admin",                                    //签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	// 如果解析失败
	if err != nil {
		return nil, err
	}

	// 如果令牌无效
	if !token.Valid {
		return nil, err
	}

	//如果令牌有效
	return mc, nil
}
