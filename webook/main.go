package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/q1ngy/Learn-Go/webook/internal/repository"
	"github.com/q1ngy/Learn-Go/webook/internal/repository/dao"
	"github.com/q1ngy/Learn-Go/webook/internal/serivce"
	"github.com/q1ngy/Learn-Go/webook/internal/web"
	"github.com/q1ngy/Learn-Go/webook/internal/web/middleware"
	//"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	server := initServer()
	initUserHandler(db, server)
	server.Run(":8080")
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := serivce.NewUserService(ur)
	uh := web.NewUserHandler(us)
	uh.RegisterRoutes(server)
}

func initServer() *gin.Engine {
	server := gin.Default()

	cors := middleware.CorsMiddlewareBuilder{}
	server.Use(cors.Build())

	login := middleware.LoginMiddlewareBuilder{}

	//store := cookie.NewStore([]byte("secret"))
	store, err := redis.NewStore(16, "tcp", "127.0.0.1:6379", "", "123456",
		[]byte("hCyJa2U47n3jrRwiLM8eXJbBUR5VJihU"),
		[]byte("F8yunGjffnhpBd6W5GyX1yMooo8ahKEp"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
	return server
}

func initDB() *gorm.DB {
	dsn := "root:123456@tcp(localhost:3306)/webook?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db.Debug()
}
