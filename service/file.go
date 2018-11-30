package service

import (
	"encoding/base64"
	"github.com/labstack/gommon/log"
	"io/ioutil"
)

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