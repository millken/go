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

func (this *MainController) Hi() {
	this.RenderHtml("<html><body>Hi</body></html>")
	log.Println("<html><body>Hi</body></html>")
}