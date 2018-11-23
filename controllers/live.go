package controllers

import (
	"github.com/astaxie/beego"
	"github.com/labstack/gommon/log"
)

//tsPath, m3u8Path, bucket

type LiveController struct {
	beego.Controller
}

// @router / [get]
func (live *LiveController) HandleLive () {
	tsPath := live.GetString("tsPath")
	m3u8Path:= live.GetString("m3u8Path")
	bucket:= live.GetString("bucket")
	log.Infof("tsPath -> %s, m3u8Path -> %s, bucket -> %s", tsPath, m3u8Path, bucket)
}


