//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:123456@tcp(localhost:3306)/webook?charset=utf8&parseTime=True&loc=Local",
	},
	Redis: RedisConfig{
		Addr:     "localhost:6379",
		Password: "1234567",
	},
}
