package controllers

import (
	"github.com/astaxie/beego"
	"github.com/labstack/gommon/log"
	"redis-agent/service"
)

type RedisController struct {
	beego.Controller
}

// @router /:key [get]
func (r *RedisController) Get() {

	key := r.GetString(":key")
	log.Info("redis receive key :", key)
	var cache = service.Cache{}
	var val string = cache.Get(key)
	log.Info("redis return val :", val)
	r.Data["json"] = val
	r.ServeJSON()
}