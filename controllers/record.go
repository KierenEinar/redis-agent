package controllers

import (
	"github.com/astaxie/beego"
	"github.com/labstack/gommon/log"
	"path"
	"redis-agent/commons"
	"redis-agent/service"
	"strings"
	"time"
)

type RecordController struct {
	beego.Controller
}

// @router / [get]
func (record *RecordController) HandleRecord (){
	tsPath := record.GetString("tsPath")
	m3u8Path:= record.GetString("m3u8Path")
	bucket:= record.GetString("bucket")
	log.Infof("tsPath -> %s, m3u8Path -> %s, bucket -> %s", tsPath, m3u8Path, bucket)

	if commons.IsBlank(&tsPath) || commons.IsBlank(&m3u8Path) || commons.IsBlank(&bucket) {
		record.Data["json"] = map[string]interface{} {"code":0, "data":"success"}
		record.ServeJSON()
		return
	}

	vod(tsPath, m3u8Path, bucket)
	record.Ctx.ResponseWriter.WriteHeader(200)
}

// @router / [post]
func (record * RecordController) OpenRecord () {
	log.Info("request body -> ", record.GetString("name"))
	vodKey := record.GetString("name")
	vodName := time.Now().Nanosecond()
	cache := service.Cache{}

	log.Info("vodKey -> ", vodKey, "," , "vodName -> ", vodName)
	cache.Set(vodKey, vodName, -1)
	log.Info("vod 写入redis 成功" )
	record.Ctx.ResponseWriter.WriteHeader(200)
}


func vod (tsPath string, m3u8Path string, bucket string) {
	log.Infof("vod,  tsKey -> %s, m3u8Key -> %s, bucket -> %s", tsPath, m3u8Path, bucket)
	data := make (chan string)
	write := make(chan string)
	go service.ReadFile(&m3u8Path, data, false)
	replace := "#EXT-X-DISCONTINUITY"
	end := "\r\n#EXT-X-ENDLIST"
	f := <- data
	f = strings.Replace(f, replace,"", -1)
	f = f + end
	temp :=m3u8Path + ".vod"
	go service.WriteFile (&temp, f, write, false)
	<-write
	hdfsprefix := beego.AppConfig.String("hdfsprefix")
	cache := service.Cache{}
	vodName := cache.Get(bucket)
	hdfs:= service.WebHdfsClient{}
	remoteFile := path.Join(hdfsprefix, vodName, "index.m3u8")
	hdfs.UploadFile(temp, remoteFile)
}