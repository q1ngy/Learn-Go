package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string

	ErrKeyNotExist       = errors.New("验证码不存在")
	ErrCodeSendTooMany   = errors.New("发送太频繁")
	ErrCodeVerifyTooMany = errors.New("验证太频繁")
)

type CodeCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, code string) (bool, error)
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return RedisCodeCache{
		cmd: cmd,
	}
}

func (c RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

func (c RedisCodeCache) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case -2:
		return false, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		return true, nil
	}
}

func (c RedisCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type LocalCodeCache struct {
	cache *cache.Cache
}

func NewLocalCodeCache(cache *cache.Cache) CodeCache {
	return &LocalCodeCache{
		cache: cache,
	}
}

func (l LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	var mutex sync.Mutex
	key := l.key(biz, phone)
	cnt := key + ":cnt"

	mutex.Lock()
	defer mutex.Unlock()

	_, left, found := l.cache.GetWithExpiration(key)
	if found {
		if left.IsZero() {
			return errors.New("验证码存在，但是没有过期时间")
		}
		if time.Until(left) > 540 {
			return ErrCodeSendTooMany
		}
	}

	l.cache.Set(key, code, 600*time.Second)
	l.cache.Set(cnt, 3, 600*time.Second)

	return nil
}

func (l LocalCodeCache) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	var mutex sync.Mutex
	key := l.key(biz, phone)
	cnt := key + ":cnt"

	mutex.Lock()
	defer mutex.Unlock()

	cntVal, found := l.cache.Get(cnt)
	if found {
		if cntVal.(int) <= 0 {
			return false, ErrCodeVerifyTooMany
		}
	} else {
		return false, nil
	}

	val, left, found := l.cache.GetWithExpiration(key)
	if found {
		if val == code {
			l.cache.Set(cnt, 0, time.Until(left))
			return true, nil
		} else {
			l.cache.Set(cnt, cntVal.(int)-1, time.Until(left))
			return false, nil
		}
	}

	return false, nil
}

func (l LocalCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
