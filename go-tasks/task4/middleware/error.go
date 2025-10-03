package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/q1ngy/go-tasks/handler"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		var errMsg string
		switch v := recovered.(type) {
		case string:
			c.JSON(http.StatusBadRequest, handler.Resp{
				Success: false,
				Message: errMsg,
			})
		case error:
			errMsg = v.Error()
			c.JSON(http.StatusInternalServerError, handler.Resp{
				Success: false,
				Message: errMsg,
			})
		default:
		}

		c.Abort()
	})
}
