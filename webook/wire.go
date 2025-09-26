//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/q1ngy/Learn-Go/webook/internal/repository"
	"github.com/q1ngy/Learn-Go/webook/internal/repository/cache"
	"github.com/q1ngy/Learn-Go/webook/internal/repository/dao"
	"github.com/q1ngy/Learn-Go/webook/internal/service"
	"github.com/q1ngy/Learn-Go/webook/internal/web"
	"github.com/q1ngy/Learn-Go/webook/ioc"
)
import "github.com/google/wire"

func InitWebServer() *gin.Engine {
	wire.Build(ioc.InitDB, ioc.InitRedis, ioc.InitCache,
		dao.NewUserDao,

		cache.NewUserCache,
		//cache.NewRedisCodeCache,
		cache.NewLocalCodeCache,

		repository.NewCachedUserRepository,
		repository.NewCodeRepository,

		ioc.InitSMSService,
		service.NewCodeService,
		service.NewUserService,

		web.NewUserHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
