// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/q1ngy/Learn-Go/mall/user/api/internal/config"
	"github.com/q1ngy/Learn-Go/mall/user/api/internal/middleware"
	"github.com/q1ngy/Learn-Go/mall/user/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config         config.Config
	CostMiddleware rest.Middleware
	UserModel      model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlxConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:         c,
		UserModel:      model.NewUserModel(sqlxConn, c.CacheRedis),
		CostMiddleware: middleware.NewCostMiddleware().Handle,
	}
}
