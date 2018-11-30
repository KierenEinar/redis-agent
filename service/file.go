package service

import (
	"encoding/base64"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"os"
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

func WriteFile (f *string, content string ,data chan string, base64Enc bool) {
	if base64Enc {
		content = base64.StdEncoding.EncodeToString([]byte(content))
	}
	ioutil.WriteFile(*f, []byte(content), os.ModeDevice)
	data <- ""
}

