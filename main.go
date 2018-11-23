package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
	_ "redis-agent/routers"
	"redis-agent/service"
	"strings"
	"time"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	var redisClusterNodes string = beego.AppConfig.String("redisClusterNodes")
	redisPoolSize , _ := beego.AppConfig.Int("redisPoolSize")
	redisMinIdleConns, _ := beego.AppConfig.Int("redisMinIdleConns")
	redisIdleTimeout,_ := time.ParseDuration(beego.AppConfig.String("redisIdleTimeout"))
	redisDialTimeout,_:=time.ParseDuration(beego.AppConfig.String("redisDialTimeout"))

	addrs := strings.Split(redisClusterNodes, ",")
	fmt.Println("redis cluster client init, addrs", addrs)
	var cache = service.Cache{addrs,redisPoolSize, redisMinIdleConns, redisIdleTimeout, redisDialTimeout}
	cache.Connect()

	var hdfsnamenode string = beego.AppConfig.String("hdfsnamenode")
	var hdfsuser string = beego.AppConfig.String("hdfsuser")

	webHdfsClient := &service.WebHdfsClient{hdfsnamenode, hdfsuser}
	webHdfsClient.Connect()

	log.Info("当前dev版本", random.New().String(12))
	beego.Run()

}
