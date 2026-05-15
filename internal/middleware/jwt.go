package middleware

import (
	"net/http"
	"strings"

	"wvp-pro-go/internal/config"
	"wvp-pro-go/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const JWTHeader = "Authorization"

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.GlobalConfig.UserSetting.InterfaceAuthentication {
			c.Next()
			return
		}

		// Try Authorization header first, then access-token
		tokenStr := c.GetHeader(JWTHeader)
		if tokenStr == "" {
			tokenStr = c.GetHeader("access-token")
		}
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, utils.Fail(utils.CodeUnauthorized, "请先登录"))
			c.Abort()
			return
		}

		// Support "Bearer <token>" format
		parts := strings.SplitN(tokenStr, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			tokenStr = parts[1]
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GlobalConfig.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, utils.Fail(utils.CodeUnauthorized, "登录已过期，请重新登录"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, utils.Fail(utils.CodeUnauthorized, "无效的Token"))
			c.Abort()
			return
		}

		// Store user info in context
		if userID, exists := claims["userId"]; exists {
			c.Set("userId", userID)
		}
		if username, exists := claims["username"]; exists {
			c.Set("username", username)
		}

		c.Next()
	}
}
