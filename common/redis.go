package common

import (
	"fmt"

	"github.com/go-redis/redis"
	"share.ac.cn/config"
)

var client *redis.Client

func InitRedisClient() {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Conf.Redis.Address, config.Conf.Redis.Port),
		Password:     config.Conf.Redis.Password,
		DB:           config.Conf.Redis.Db,
		PoolSize:     config.Conf.Redis.Size,
		MinIdleConns: config.Conf.Redis.ConnMax,
	})
}

func GetClient() *redis.Client {
	return client
}
