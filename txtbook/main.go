package main

import (
	"github.com/qiniu/log"
	//"os"
	//"test"
	//"config"
	//"Book"
	"web"
)


func main() {

	//w, _ := os.OpenFile("debug", os.O_CREATE|os.O_APPEND, 0644)
	//log.SetOutput(w)
	//log.SetOutputLevel(log.Ldebug)

	log.Debugf("Server starting...\n")
	//log.Debug(config.GetRedisConfig().Addr)
	//test.Testlog()
	//Book.AddUrl("http://aaa.com/sdf","sd","sd")
	web.Run()
}
