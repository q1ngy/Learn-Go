package ioc

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func InitCache() *cache.Cache {
	return cache.New(10*time.Minute, 10*time.Minute)
}
