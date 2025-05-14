package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"awesomeProject/utils"
)

// 日志中间件
func SetLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 记录请求基本信息
		utils.Logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"referer":    c.Request.Referer(),
			"protocol":   c.Request.Proto,
		}).Info("Request received")

		// 如果有请求体，可以记录部分内容（注意不要记录敏感信息）
		if c.Request.ContentLength > 0 && c.Request.ContentLength < 1024 {
			body, err := io.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // 重新设置body供后续使用
				utils.Logger.WithField("body", string(body)).Debug("Request body")
			}
		}

		// 继续处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(startTime)

		// 记录响应信息
		fields := logrus.Fields{
			"status":    c.Writer.Status(),
			"latency":   latency.String(),
			"bytes":     c.Writer.Size(),
			"client_ip": c.ClientIP(),
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
		}

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
			utils.Logger.WithFields(fields).Error("Request completed with errors")
		} else {
			utils.Logger.WithFields(fields).Info("Request completed")
		}
	}
}
