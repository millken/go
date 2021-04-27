package main

import (
	"mvc"
	"github.com/millken/go-whois"
	"logger"
    "regexp"
)
var (
    regex_server = regexp.MustCompile(`[Name Server| Whois Server]: (.*)`)
)

type WhoisController struct {
	mvc.Controller
}

func (this *WhoisController) Domain() {
    var info string
    var err error
    domain := this.Request.URL.Query().Get("domain")
    server := this.Request.URL.Query().Get("server")
    if server != "" {
        info, err = whois.WhoisByServer(domain,server)
    }else{
        info, err = whois.Whois(domain)
    }
    if err != nil {
        log.Println(err)
    } else {
        ws := GetWhoisServer(info)
        this.RenderText(ws + info)
    }
}

func GetWhoisServer(body string) (result string) {
    servers := regex_server.FindAllStringSubmatch(body, -1)
    len_server := len(servers)
    if len_server > 0 {
        result = servers[len_server - 1][1]
    }
    return
}
