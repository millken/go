package main

import (
	"fmt"
	"runtime"
	"net/http"
	"io/ioutil"
	"regexp"
	"strings"
)

func main() {
	var _ = fmt.Print
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	proxyListTemplateUrl := "http://www.xici.net.co/nn/%02v"
	
	pageIndex := 1
	for {
		proxyListUrl := fmt.Sprintf(proxyListTemplateUrl, pageIndex)
		fmt.Println("# Getting", proxyListUrl)
		response, err := http.Get(proxyListUrl)
		if err != nil {
			// Most likely, all the pages have been downloaded
			fmt.Println("# ", err.Error())
			break
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("# ", err.Error())
			continue
		}
		bodyString := string(body)
		
		exp_ip, error := regexp.Compile(`<td>\d+\.\d+\.\d+\.\d+</td>`)
		exp_port, _ := regexp.Compile(`<td>([0-9]+)</td>`)
		//    exp, error := regexp.Compile(`<div class="today-results	clearfix">([0-9]+)</div>`)

		if error == nil {
			ip := exp_ip.FindAllString(bodyString, -1)
			port := exp_port.FindAllString(bodyString, -1)
			if ip != nil && len(ip) > 0 {
				for i := 0; i < len(ip); i++ {
					fmt.Printf(strings.Trim(ip[i], "</td>") + ":" + strings.Trim(port[i], "</td>") + "\n")
				}
			} else {
				fmt.Printf("#No matches found.\n")
			}

		} else {
			fmt.Printf("#", error.Error())
		}

		pageIndex++
	} 
	
}
