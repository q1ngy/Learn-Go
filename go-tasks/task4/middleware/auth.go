package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/q1ngy/go-tasks/handler"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		pattern := `^/post/([1-9]\d*)$`
		re := regexp.MustCompile(pattern)
		if path == "/user/register" || path == "/user/login" || path == "/post/all" || (re.MatchString(path) && ctx.Request.Method == "GET") {
			ctx.Next()
			return
		}

		header := ctx.Request.Header.Get("Authorization")
		segments := strings.Split(header, " ")
		if header == "" || len(segments) != 2 {
			ctx.JSON(http.StatusOK, handler.Resp{
				Success: false,
				Message: "token err",
			})
			ctx.Abort()
			return
		}
		tokenStr := segments[1]

		var claims = handler.JWTClaims{}
		_, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
			return []byte("123456"), nil
		})
		if err != nil {
			ctx.JSON(http.StatusOK, handler.Resp{
				Success: false,
				Message: "token err",
			})
		}

		ctx.Set("userId", claims.UserID)
		ctx.Next()
	}
}
