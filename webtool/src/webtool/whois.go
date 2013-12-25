package main

import (
	"mvc"
	"github.com/millken/go-whois"
	"logger"
)

type WhoisController struct {
	mvc.Controller
}

func (this *WhoisController) Domain() {
    query := this.Request.URL.Query().Get("domain")
    params := make(map[string]string)

    info, err := whois.Whois(query, params)
    if err != nil {
        log.Println(err)
    } else {
        this.RenderText(info)
    }
}
