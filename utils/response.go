package utils

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Response 基础响应结构体
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// Success 成功响应
func BuildSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func BuildErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ServerError 记录错误日志并返回 500 响应
func BuildServerError(c *gin.Context, msg string, err error) {
	zap.L().Error(msg,
		zap.Error(err),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("ip", c.ClientIP()),
	)
	c.JSON(500, Response{
		Code:    500,
		Message: msg,
		Data:    nil,
	})
}

// Custom 自定义响应
func BuildCustomResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
