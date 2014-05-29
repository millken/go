/*
 * curl http://46.16.226.10:8080/ic.asp -H "Host: iframe.ip138.com"
 * http://toop123.duapp.com/pcip/getip.php?r=0.37779340663140526
 */
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"net/url"
	"time"
	"net/http"
	"regexp"
)

var (
	nConnectFlag = flag.Int("concurrent", 100, "Number of concurrent connections")
	dataFileFlag = flag.String("data", "", "proxy file")
	urlFlag      = flag.String("url", "http://iframe.ip138.com/ic.asp", "checkurl, default is : http://iframe.ip138.com/ic.asp")
	ipFlag       = flag.String("ip", "", "proxy ip address example 8.8.8.8:88")
	outputFlag   = flag.String("output", "", "output file")
)

func init() {
	flag.Parse()
}

/*
//https://github.com/mattn/hugo/tree/1e7f18e04e5baf1975da6b87403fa66afe057535/commands
func tweakLimit() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Unable to obtain rLimit", err)
	}
	if rLimit.Cur < rLimit.Max {
		rLimit.Max = 999999
		rLimit.Cur = 999999
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			fmt.Println("Unable to increase number of open files limit", err)
		}
	}
}
*/
func checkproxy(ip string) {
	fmt.Println("Checking:", ip)
	proxyUrl, err := url.Parse("http://" + ip)
	_, err = net.DialTimeout("tcp", ip, 5*time.Second)
	if err != nil {
		fmt.Printf("%s => %s\n", ip, err.Error())
		return		
	}
	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	response, err := httpClient.Get(*urlFlag)
	if err != nil {
		fmt.Printf("%s => %s\n", ip, err.Error())
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s => %s\n", ip, err.Error())
		return
	}

	bodyString := string(body)

	if strings.Index(bodyString, "<body") < 0 && strings.Index(strings.ToUpper(bodyString), "IP") < 0 {
		fmt.Printf("%s => Reveived page is not HTML\n", ip)
		return
	}

	if *outputFlag != "" {
		if err := WriteFile(*outputFlag, 0600, fmt.Sprintf("%s\n", ip)); err != nil {
			fmt.Printf("write file error : %s", err)
		}
	}
	fmt.Printf("[%s] ->%s : ok\n", ip, *urlFlag)
}

func worker(lChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for ip := range lChan {
		checkproxy(ip)
		//fmt.Printf("Done processing ip #%s\n", ip)
	}

}
func WriteFile(fname string, perm uint32, str string) error {

	var file *os.File

	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))

	if err != nil {
		return err
	}

	if _, err := file.WriteString(str); err != nil {
		return err
	}
	file.Close()
	return err
}
func main() {

	lCh := make(chan string)
	wg := new(sync.WaitGroup)

	// Adding routines to workgroup and running then
	for i := 0; i < *nConnectFlag; i++ {
		wg.Add(1)
		go worker(lCh, wg)
	}

	if *ipFlag != "" {
		lCh <- *ipFlag
	}

	if *dataFileFlag != "" {
		file, err := os.Open(*dataFileFlag) // For read access.
		defer file.Close()
		if err != nil {
			fmt.Printf("open file %s error: %v\n", *dataFileFlag, err)
			return
		}

		data, err := ioutil.ReadAll(file)
		iplist := strings.Split(string(data), "\n")
		for _, ips := range iplist {
			ipreg := regexp.MustCompile(`\d+.\d+.\d+.\d+:\d+`)
			newip := ipreg.FindStringSubmatch(ips)
			if newip == nil || len(newip) < 1 || len(newip[0]) < 10 {
				continue
			}
			//fmt.Printf("%s\n", newip[0])
			ip := strings.Trim(newip[0], "\r\n ")
			lCh <- ip

		}
	}
	close(lCh)

	wg.Wait()
}
