package main

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "redis-agent/routers"
	"redis-agent/service"
	"strings"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	var redisClusterNodes string = beego.AppConfig.String("redisClusterNodes")
	addrs := strings.Split(redisClusterNodes, ",")
	fmt.Println("redis cluster client init, addrs", addrs)
	var cache = service.Cache{addrs}
	cache.Connect()

	var hdfsnamenode string = beego.AppConfig.String("hdfsnamenode")
	var hdfsuser string = beego.AppConfig.String("hdfsuser")

	webHdfsClient := &service.WebHdfsClient{hdfsnamenode, hdfsuser}
	webHdfsClient.Connect()

	beego.Run()
}
