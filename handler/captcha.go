package handler

import (
	"fmt"
	"net/http"

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate captcha",
		})
		return
	}

	utils.BuildSuccessResponse(c, gin.H{
		"success":      true,
		"captchaID":    captchaID,
		"captchaImage": fmt.Sprintf("data:image/png;base64,%s", captchaBase64),
	})
}

// 验证验证码
func VerifyCaptcha(c *gin.Context) {
	// 从请求中获取 captchaID 和用户输入的验证码
	type VerifyRequest struct {
		CaptchaID    string `json:"captchaID" binding:"required"`
		CaptchaValue string `json:"captchaValue" binding:"required"`
	}

	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
		})
		return
	}

	// 验证验证码（会同时检查是否过期）
	isValid := utils.CaptchaStore.Verify(req.CaptchaID, req.CaptchaValue, true)
	if isValid {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Captcha verified successfully",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Captcha verification failed or expired",
		})
	}
}
