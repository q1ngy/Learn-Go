package test_once

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var once sync.Once
var wg sync.WaitGroup
var client *Client

type Client struct {
}

func (c *Client) initClient() *Client {
	once.Do(func() {
		client = &Client{} // 可能很耗时的构建
	})
	return client
}

func TestOnce(t *testing.T) {
	//wg.Add(2)
	go once.Do(initConfig)
	go once.Do(initConfig)

	time.Sleep(5 * time.Second)

	//wg.Wait()
}

func initConfig() {
	//wg.Done()
	fmt.Println("初始化配置...")
}
