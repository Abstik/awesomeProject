package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"awesomeProject/dao"
	"awesomeProject/utils"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, "用户需要登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, "无效的token")
			c.Abort()
			return
		}

		mc, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, "无效的token")
			c.Abort()
			return
		}

		// 校验用户是否仍然存在
		if _, err := dao.GetMemberByUsername(mc.UserName); err != nil {
			c.JSON(http.StatusUnauthorized, "用户不存在或已被删除")
			c.Abort()
			return
		}

		c.Set("userID", mc.UserID)
		c.Set("username", mc.UserName)
		c.Set("status", mc.Status)
		c.Next()
	}
}

func IsAdminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status, exists := ctx.Get("status")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "用户需要登录"})
			ctx.Abort()
			return
		}
		statusInt, ok := status.(int)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "status 类型异常"})
			ctx.Abort()
			return
		}
		// 只有管理员(角色id为0)才有权限
		if statusInt != 0 {
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "无管理员权限"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
