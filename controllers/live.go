package controllers

import (
	"github.com/astaxie/beego"
	"github.com/labstack/gommon/log"
	"path"
	"redis-agent/commons"
	"redis-agent/service"
)

//tsPath, m3u8Path, bucket

type LiveController struct {
	beego.Controller
}

// @router / [get]
func (live *LiveController) HandleLive (){
	tsPath := live.GetString("tsPath")
	m3u8Path:= live.GetString("m3u8Path")
	bucket:= live.GetString("bucket")
	log.Infof("tsPath -> %s, m3u8Path -> %s, bucket -> %s", tsPath, m3u8Path, bucket)

	if commons.IsBlank(&tsPath) || commons.IsBlank(&m3u8Path) || commons.IsBlank(&bucket) {
		live.Data["json"] = map[string]interface{} {"code":0, "data":"success"}
		live.ServeJSON()
		return
	}

	Live(tsPath, m3u8Path, bucket)

	live.Ctx.ResponseWriter.WriteHeader(200)

}

func Live (tsPath string, m3u8Path string, bucket string) {

	tsFile := make(chan string, 1)
	m3u8File := make (chan string, 1)

	go service.ReadFile(&tsPath,tsFile, true)
	go service.ReadFile(&m3u8Path,m3u8File,false)


	//写入ts文件到redis, key为 推流码 + 文件名
	tsKey:= bucket + "/" + path.Base(tsPath)
	m3u8Key:= bucket + "/" + path.Base(m3u8Path)

	log.Infof("tsKey -> %s, m3u8Key -> %s", tsKey, m3u8Key)

	tsRaw := <-tsFile
	m3u8Raw := <-m3u8File
	redis := service.Cache{}
	redis.Set(m3u8Key, m3u8Raw, 60)
	redis.Set(tsKey, tsRaw,60)

	log.Info("写入 redis 成功")
}

