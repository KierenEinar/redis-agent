package service

import (
	"github.com/go-redis/redis"
	"github.com/labstack/gommon/log"
	"os"
	"time"
)

var redisClient *redis.ClusterClient

type Cache struct{
	Addrs []string
	PoolSize int
	MinIdleConns int
	IdleTimeout time.Duration
	DialTimeout time.Duration
}

func (c *Cache) Connect () {
	rediClusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: c.Addrs,
		DialTimeout: c.DialTimeout,
		PoolSize: c.PoolSize,
		MinIdleConns:c.MinIdleConns,
		IdleTimeout: c.IdleTimeout,
	})

	err := rediClusterClient.Ping().Err()
	if err != nil {
		log.Error("connect redis cluster error", err)
		os.Exit(-1)
	} else {
		log.Info("connect redis cluster success")
	}

	redisClient = rediClusterClient

}

func (*Cache) Set(key string, value interface{}, seconds int64) {
	redisClient.Set(key, value, time.Duration(seconds) * time.Second)
}

func (*Cache) Get(key string) (string){
	 result := redisClient.Get(key)
	 return result.Val()
}