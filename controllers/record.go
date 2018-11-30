package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/labstack/gommon/log"
	"redis-agent/commons"
	"redis-agent/models"
	"redis-agent/service"
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
	var obj models.Record
	json.Unmarshal(record.Ctx.Input.RequestBody, &obj)
	log.Info("open record, name -> ", obj.Name)
	vodName := obj.Name
	vodKey := time.Now().Nanosecond()
	cache := service.Cache{}
	log.Infof("vodName -> %s,	vodKey -> %s ", vodName, vodKey)
	cache.Set(vodName, vodKey, 60)
	record.Ctx.ResponseWriter.WriteHeader(200)
}


func vod (tsPath string, m3u8Path string, bucket string) {
	log.Infof("vod,  tsKey -> %s, m3u8Key -> %s, bucket -> %s", tsPath, m3u8Path, bucket)
}