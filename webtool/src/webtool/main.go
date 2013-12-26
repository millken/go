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
		go sigHandler()
}

func main() {
	App.ServeListen("0.0.0.0", 82)
	App.ServeFile("/favicon.ico", "../static/favicon.ico")
	App.ServeFile("/sitemap.xml", "../static/sitemap.xml")
	App.AddPreAction(&MainController{}, "Hi")
    App.Router.AddRoute("127.0.0.1", "/hello/:user", &MainController{}, "Hello2")
    App.Router.AddRoute("*", "/whois/:domain", &WhoisController{}, "Domain")
	App.Router.AddStaticPath("*", "/static/", "../static/")
    App.AddTemplates("../template/default/")
 	App.Run()
}
