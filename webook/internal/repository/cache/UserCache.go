package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/q1ngy/Learn-Go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (uc *UserCache) Get(ctx *gin.Context, uid int64) (domain.User, error) {
	key := uc.key(uid)
	result, err := uc.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(result), &u)
	return u, err
}

func (uc *UserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func (uc *UserCache) Set(ctx *gin.Context, du domain.User) error {
	key := uc.key(du.Id)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return uc.cmd.Set(ctx, key, data, uc.expiration).Err()
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
