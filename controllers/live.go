package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/labstack/gommon/log"
	"io/ioutil"
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


	if commons.IsBlank(&tsPath) || commons.IsNotBlank(&m3u8Path) || commons.IsNotBlank(&bucket) {
		live.Data["json"] = map[string]interface{} {"code":0, "data":"success"}
		live.ServeJSON()
		return
	}

	tsFile := make(chan string, 1)
	m3u8File := make (chan string, 1)

	go ReadFile(&tsPath,tsFile, true)
	go ReadFile(&m3u8Path,m3u8File,false)


	//写入ts文件到redis, key为 推流码 + 文件名
	tsKey:= bucket + "/" + path.Base(tsPath)
	m3u8Key:= bucket + "/" + path.Base(m3u8Path)

	log.Infof("tsKey -> %s, m3u8Key -> %s", tsKey, m3u8Key)


	tsRaw := <-tsFile
	m3u8Raw := <-m3u8File

	log.Infof("tsRaw -> %s, m3u8Raw -> %s", tsRaw, m3u8Raw)

	redis := service.Cache{}
	redis.Set(m3u8Key, m3u8Raw, 60)
	redis.Set(tsKey, tsRaw,60)

	live.Data["json"] = map[string]interface{} {"code":0, "data":"write redis success"}
	live.ServeJSON()
}


func ReadFile (f *string, data chan string, base64Enc bool) {
	raw,err:= ioutil.ReadFile(*f)
	if err!=nil {
		log.Error("读取文件失败, 文件->", *f)
	}
	var str string
	if base64Enc {
		str = base64.StdEncoding.EncodeToString(raw)
	} else {
		str= string(raw[:])
	}
	data <- str
}