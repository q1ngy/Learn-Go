package ioc

import (
	"github.com/q1ngy/Learn-Go/webook/internal/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
	})
	return client
}
