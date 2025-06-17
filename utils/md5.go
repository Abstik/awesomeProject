package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
)

// 对密码进行加密的盐
const secret = "3gmobile"

// 对密码进行加密
func EncryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// 生成5位随机数字字符串
func Rand5Digits() string {
	return fmt.Sprintf("%05d", rand.Intn(100000))
}
