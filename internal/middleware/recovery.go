package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"wvp-pro-go/internal/utils"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic recovered", zap.Any("error", err))
				c.JSON(http.StatusOK, utils.Fail(utils.CodeServerError, "系统异常"))
				c.Abort()
			}
		}()
		c.Next()
	}
}
