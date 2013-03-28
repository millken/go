package test

import (
	"github.com/qiniu/log"
	"time"
)

func Testlog() {

	for {
		time.Sleep(time.Second)
		log.Debug("Debug in log.go")
	}
}