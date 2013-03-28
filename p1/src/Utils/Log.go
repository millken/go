//https://github.com/sun8588/boboqueue
package Utils

import (
	"log"
	"os"
)

var LogOut *log.Logger
var files os.File

type logMsg struct {
	format string
	value  []interface{}
}

var logChan chan *logMsg

func InitLogOut(logFile string,logChanNum int) error {
	//设置log文件
	files, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		return LogErr(err)
	}
	LogOut = log.New(files, "", 0)
	LogOut.SetFlags(log.Ldate | log.Ltime)
	//init log chan
	logChan = make(chan *logMsg, logChanNum)
	go writeLog()
	return nil
}
func DeferFiles() {
	files.Close()
	close(logChan)
}
func LogInfo(format string, info ...interface{}) {
	if logChan != nil {
		logChan <- &logMsg{format, info}
//		logChan <- &logMsg{"now logchan=%d\n",[]interface{}{len(logChan)}}
	} else {
		log.Printf(format, info...)
	}
}

func writeLog() {
	var logInfo *logMsg
	for {
		select {
		case logInfo = <-logChan:
			if LogOut != nil {
				LogOut.Printf((*logInfo).format, (*logInfo).value...)
			} else {
				log.Printf((*logInfo).format, (*logInfo).value...)
			}
		}
	}
}
