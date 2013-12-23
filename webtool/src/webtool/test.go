package main

import (
	"mvc"
	"logger"
)

type MainController struct {
	mvc.Controller
}

func (this *MainController) Hello() {
	this.RenderText("hello world")
	log.Println("hello world")
}

func (this *MainController) Hello2() {
	this.RenderText("hello world 2")
}

func (this *MainController) Hi() {
	//this.RenderHtml("<html><body>Hi</body></html>")
	log.Debugf("Host:%s, Path:%s, RawQuery: %s", this.Request.Host, this.Request.URL.Path, this.Request.URL.RawQuery)
	log.Println("<html><body>Hi</body></html>")
}