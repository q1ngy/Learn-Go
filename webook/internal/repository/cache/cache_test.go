package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
)

func TestCache(t *testing.T) {
	c := cache.New(-1, 0)
	c.Set("a", 42, 3*time.Second)
	v, b := c.Get("a")
	fmt.Println(v, b)
	expiration, left, b2 := c.GetWithExpiration("a")
	fmt.Println(expiration, left, b2)

	until := time.Until(left)
	fmt.Println(until)

}
