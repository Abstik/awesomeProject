package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"

	"awesomeProject/utils"
)

// 生成验证码
func GetCaptcha(c *gin.Context) {
	// 配置验证码：数字类型验证码
	driver := base64Captcha.NewDriverDigit(
		80,  // 高度
		240, // 宽度
		5,   // 验证码长度
		0.7, // 噪点强度
		80,  // 字体大小
	)

	// 创建验证码实例
	captcha := base64Captcha.NewCaptcha(driver, utils.CaptchaStore)

	// 生成验证码
	captchaID, captchaBase64, _, err := captcha.Generate()
	if err != nil {
		utils.BuildErrorResponse(c, 500, "生成验证码失败")
		return
	}

	utils.BuildSuccessResponse(c, gin.H{
		"captchaID":    captchaID,
		"captchaImage": captchaBase64,
	})
}

// 验证验证码
func VerifyCaptcha(CaptchaID, CaptchaValue string) bool {
	// 验证验证码（检查是否正确和是否过期）
	isValid := utils.CaptchaStore.Verify(CaptchaID, CaptchaValue, true)
	if !isValid {
		return false
	}
	return true
}
