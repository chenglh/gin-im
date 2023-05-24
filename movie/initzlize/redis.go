package initzlize

import (
	"IM/movie/global"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

// 连接Redis数据库
func InitRedis() (err error) {
	global.Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password: "",
		DB:       0,
	})

	if _, err := global.Rdb.Ping(context.Background()).Result(); err != nil {
		return err
	}

	return nil
}
