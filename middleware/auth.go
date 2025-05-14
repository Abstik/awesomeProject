package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"awesomeProject/utils"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URL
		// 这里Token放在请求头Header的Authorization中，并使用Bearer开头
		// 格式：Authorization: Bearer xxx.xxx.xxx
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		//如果未携带令牌，代表没有登录
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, "用户需要登录")
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, "无效的token")
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, "无效的token")
			c.Abort()
			return
		}

		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set("userID", mc.UserID)
		c.Set("username", mc.UserName)
		c.Set("status", mc.Status)
		c.Next() // 后续的处理函数可以用过c.Get(ContextUserIDKey)来获取当前请求的用户信息
	}
}

func IsAdminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status, _ := ctx.Get("status")
		statusInt := status.(int)
		// 只有管理员(角色id为0)才有权限
		if statusInt != 0 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无管理员权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}

}
