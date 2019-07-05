package redisHelper

import (
	"sync"

	"github.com/go-redis/redis"
)

var instance *redis.Client
var once sync.Once
var conf *redisConf

type redisConf struct {
	addr     string
	password string
	db       int
}

func init() {
	_ = &redisConf{"localhost:6379", "", 0}
}

func GetRedisClient() *redis.Client {
	once.Do(func() {
		instance := redis.NewClient(&redis.Options{
			Addr:     conf.addr,
			Password: conf.password,
			DB:       conf.db,
		})
		_, err := instance.Ping().Result()
		if err != nil {

		}
	})
	return instance
}
