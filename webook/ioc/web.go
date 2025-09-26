package ioc

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/q1ngy/Learn-Go/webook/internal/web"
	"github.com/q1ngy/Learn-Go/webook/internal/web/middleware"
	"github.com/q1ngy/Learn-Go/webook/pkg/ginx/middleware/ratelimit"
	"github.com/redis/go-redis/v9"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(cmd redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.CorsBuild(),
		ratelimit.NewBuilder(cmd, time.Second, 100).Build(),
		(&middleware.LoginJWTMiddlewareBuilder{}).CheckLogin(),
	}
}

func useSession(server *gin.Engine) {
	login := middleware.LoginMiddlewareBuilder{}

	store := cookie.NewStore([]byte("secret"))
	//store, err := redis.NewStore(16, "tcp", "127.0.0.1:6379", "", "123456",
	//	[]byte("hCyJa2U47n3jrRwiLM8eXJbBUR5VJihU"),
	//	[]byte("F8yunGjffnhpBd6W5GyX1yMooo8ahKEp"))
	//if err != nil {
	//	panic(err)
	//}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
