package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"github.com/labstack/gommon/log"
	"os"
	"strings"
)

func init () {

	var redisClusterNodes string = beego.AppConfig.String("redisClusterNodes")

	addrs := strings.Split(redisClusterNodes, ",")

	fmt.Println("redis cluster client init, addrs", addrs)

	rediClusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
		DialTimeout: 500000000,
	})

	err := rediClusterClient.Ping().Err()
	if err != nil {
		log.Error("connect redis cluster error", err)
		os.Exit(-1)
	} else {
		log.Info("connect redis cluster success")
	}
}
