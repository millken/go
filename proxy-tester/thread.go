package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"net/http"
	"net/url"
	"runtime"
	"net"
	"time"
	"io"
	"os"
	"os/signal"
)
const (
	// How many requests/responses can be in the queue before blocking.
	QUEUE_LENGTH = 65535
)
type Wip struct {
	Ip string
}
var global_wip = make(chan *Wip, QUEUE_LENGTH)
var global_file, _ = os.OpenFile("checkedproxy.txt", os.O_APPEND, 0666) 

func main() {
	var time1 = time.Now().Unix()
	runtime.GOMAXPROCS(2)
	fmt.Printf("Started: %d cores, %d threads", 2, 64)

	for i := 0; i < 128; i++ {
		go start_worker(global_wip)
	}
	content, err := ioutil.ReadFile("xici2.list")
	if err != nil {
		panic(err.Error())
	}
	lines := strings.Split(string(content), "\n")

	for i := 0; i < len(lines); i++ {
		var line = strings.Trim(lines[i], " \t\n\r")
		if len(line) == 0 { continue }
		if line[0] == '#' { continue }
		var tokens = strings.Split(line, " ")
		if len(tokens) <= 0 { continue }
		var proxy = strings.Trim(tokens[0], " \t\n\r")
		if len(proxy) < 7 { continue }

		global_wip <- &Wip{proxy}
	}
	terminate := make(chan os.Signal)
	signal.Notify(terminate, os.Interrupt)

	<-terminate
	global_file.Close()
	fmt.Printf("total time: %d seconds\n", time.Now().Unix() - time1)
	fmt.Printf(" signal received, stopping")
}

func checkIp(proxy string)(success bool, errorMessage string) {
	fmt.Println("Checking:", proxy)
	proxyUrl, err := url.Parse("http://" + proxy)
	_, err = net.DialTimeout("tcp", proxy,  5 * time.Second)
	if err != nil {
		return false, err.Error()
	}
	httpClient := &http.Client {Transport: &http.Transport { Proxy: http.ProxyURL(proxyUrl) } }
	response, err := httpClient.Get("http://iframe.ip138.com/ic.asp")
	if err != nil { return false, err.Error() }

	body, err := ioutil.ReadAll(response.Body)
	if err != nil { return false, err.Error() }

	bodyString := string(body)

	if strings.Index(bodyString, "<body") < 0 && strings.Index(bodyString, "<head") < 0 {
		return false, "Reveived page is not HTML"
	}
	return true, bodyString
}

func start_worker(queue <-chan *Wip) {
	for w := range queue {
		status, err := checkIp(w.Ip)
		if status {
			io.WriteString(global_file, w.Ip + "\n") 
			fmt.Printf("OK %s", w.Ip)
		}else{
			fmt.Printf("#", err)
		}
		
	}
}
