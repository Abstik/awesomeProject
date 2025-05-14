package utils

import (
	"time"

	"github.com/mojocn/base64Captcha"
)

// 全局存储器，设置过期时间为 1 分钟
var CaptchaStore = base64Captcha.NewMemoryStore(10240, 1*time.Minute) // 最大存储 10240 个验证码，过期时间为 1 分钟
