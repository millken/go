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
	log.Debugf("Host:%s, Path:%s, RawQuery: %s", this.Request.Host, this.Request.URL.Path, this.Request.URL.RawQuery)
	this.RenderHtml("<html><body>hello, " + this.Request.URL.Query().Get("user")+ "</body></html>")
}

func (this *MainController) Hi() {
	log.Debugf("Host:%s, Path:%s, RawQuery: %s", this.Request.Host, this.Request.URL.Path, this.Request.URL.RawQuery)
}