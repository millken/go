package main

import (
	"mvc"
	"logger"
)

type MainController struct {
	mvc.Controller
}

func (this *MainController) Hello() {
	log.Println("hello world")
}

func (this *MainController) Hi() {
	log.Println("<html><body>Hi</body></html>")
}