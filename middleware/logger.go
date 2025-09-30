package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"awesomeProject/utils"
)

func ZapLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.EnsureZapLoggerReady()

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next() // 执行后续中间件

		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()

		utils.Logger.Info("",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		)
	}
}
