package main

import (
	"mvc"
	//"logger"
)

type HomeController struct {
	mvc.Controller
}

func (this *HomeController) Index() {
	this.Display("home.html")
}
