package redisHelper

import (
	"sync"

	"github.com/go-redis/redis"
)

var instance *redis.Client
var once sync.Once
var Conf *redisConf

type RedisConf struct {
	addr     string
	password string
	db       int
}

func init() {
}

func GetRedisClient() *redis.Client {
	once.Do(func() {
		instance = redis.NewClient(&redis.Options{
			Addr:     Conf.addr,
			Password: Conf.password,
			DB:       Conf.db,
		})
		_, err := instance.Ping().Result()
		if err != nil {

		}
	})
	return instance
}
