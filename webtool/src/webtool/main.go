package main

import (
	"os"
	"os/signal"
	"logger"
	"mvc"
)

var VERSION string = "1.0"
var gitVersion string
var (
        App = mvc.NewApp()
)
func sigHandler() {
        terminate := make(chan os.Signal)
        signal.Notify(terminate, os.Interrupt)

        <-terminate
        log.Printf("signal received, stopping")
		os.Exit(0)
}

func init() {
        if len(gitVersion) > 0 {
                VERSION = VERSION + "/" + gitVersion
        }

        log.SetOutputLevel(log.Ldebug)
}

func main() {
	go sigHandler()
    log.Println("start server")
	App.AddPreAction(&MainController{}, "Hi")
	App.Run()
}
