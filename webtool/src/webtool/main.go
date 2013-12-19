package main

import (
	"os"
	"os/signal"
	"logger"
)

var VERSION string = "1.0"
var gitVersion string

func sigHandler() {
        terminate := make(chan os.Signal)
        signal.Notify(terminate, os.Interrupt)

        <-terminate
        log.Printf("signal received, stopping")
}

func init() {
        if len(gitVersion) > 0 {
                VERSION = VERSION + "/" + gitVersion
        }

        log.SetOutputLevel(log.Ldebug)
}

func main() {
	go sigHandler()
}
