package main

import (
    "fmt"
    "net/http"
    "runtime"
)

var urllist = [...]string{
    "http://www.baidu.com",
    "http://www.google.com",
    "http://www.360.com",
    "http://www.qq.com",
    "http://www.sina.com",
    "http://www.sohu.com",
    "http://www.jd.com",
    "http://www.taobao.com",
    "http://www.chingyu.net",
    "http://www.dajie.com",
    "http://www.56.com",
    "http://www.163.com",
    "http://www.weibo.com",
    "http://www.263.com",
    "http://www.duowan.com",
    "http://www.duokan.com",
    "http://www.youku.com",
    "http://www.tudou.com",
    "http://www.zhihu.com",
    "http://www.21cn.com",
    "http://www.chinaunix.com",
    "http://www.oschina.com",
    "http://www.iteye.com",
}

var worker = runtime.NumCPU()

func gethtml(url string, respChan chan []string) {
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
        return
    }
    respChan <- resp.Header["Server"]
    resp.Body.Close()
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    respChan := make(chan []string)

    for _, url := range urllist {
        go gethtml(url, respChan)
    }

    for i := 0; i < len(urllist); i++ {
        fmt.Println(<-respChan)
    }
}
