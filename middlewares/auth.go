package middlewares

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定

		// access token 为空的情况
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseAccessToken(parts[1])
		if err != nil {
			// 如果access token 过期，通知前端给一下refresh token
			if err == jwt.ErrorExpiredAccessToken {
				refreshToken, _ := c.Cookie("refresh_token")
				rc, err := jwt.ParseRefreshToken(refreshToken)
				if err != nil {
					controller.ResponseError(c, controller.CodeNeedLogin)
					c.Abort()
					return
				}
				if rc.Username == mc.Username {
					newRefreshToken, _ := jwt.GenRefreshToken(mc.UserID, mc.Username, mc.UserType)
					newAccessToken, _ := jwt.GenAccessToken(mc.UserID, mc.Username, mc.UserType)
					//更新refresh_token
					c.SetCookie("refresh_token",
						newRefreshToken,
						viper.GetInt("auth.refresh_token_expire"),
						"/",
						"",
						false,
						true)

					// 将新的Access token 通过相应头返回给前端
					c.Writer.Header().Set("X-New-Access-Token", newAccessToken)
					// 更新请求中的Authorization 头，继续处理原请求
					c.Request.Header.Set("Authorization", "Bearer "+newAccessToken)

				}
				// controller.ResponseError(c, controller.CodeExpiredAccessToken)
				// c.Abort()
				// return
			} else {
				controller.ResponseError(c, controller.CodeInvalidToken)
				c.Abort()
				return
			}

		}

		// rc, err := jwt.ParseRefreshToken(refreshToken.Token)
		// if err != nil {
		// 	if err == jwt.ErrorExpiredRefreshToken {
		// 		controller.ResponseError(c, controller.CodeNeedLogin)
		// 		c.Abort()
		// 		return
		// 	}

		// 	controller.ResponseError(c, controller.CodeInvalidToken)
		// 	c.Abort()
		// 	return
		// } else {
		// 	if rc.Username == mc.Username {
		// 		accessToken, _ := jwt.GenAccessToken(rc.UserID, rc.Username, "normal")
		// 		refreshToken, _ := jwt.GenRefreshToken(rc.UserID, rc.Username, "normal")
		// 		c.JSON(http.StatusOK, gin.H{
		// 			"access_token":  accessToken,
		// 			"refresh_token": refreshToken,
		// 		})
		// 		c.Abort()
		// 		return
		// 	}
		// }
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("userID")来获取当前请求的用户信息
	}
}
