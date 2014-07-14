package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var (
	nConnectFlag = flag.Int("concurrent", 100, "Number of concurrent connections")
	dataFileFlag = flag.String("data", "", "File containing message to send to responsive hosts ('%s' will be replaced with host IP)")
	urlFlag   = flag.String("url", "http://www.google.com", "checkurl, default is : http://www.google.com")
	wordFlag = flag.String("word", "<title>", "match word")
	ipFlag       = flag.String("ip", "", "ip address(es) cidr to scan. example 8.8.8.8/24")
	outputFlag   = flag.String("output", "", "output file")
)

func init() {
	flag.Parse()
	// Increase file descriptor limit
	rlimit := syscall.Rlimit{Max: uint64(*nConnectFlag + 50000), Cur: uint64(*nConnectFlag + 50000)}
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting rlimit: %s", err)
	}
}

//http://play.golang.org/p/TZbIBev4pU

func ipStringToI32(a string) uint32 {
	return ipToI32(net.ParseIP(a))
}
func ipToI32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func htons(port uint16) uint16 {
	var (
		lowbyte  uint8  = uint8(port)
		highbyte uint8  = uint8(port << 8)
		ret      uint16 = uint16(lowbyte)<<8 + uint16(highbyte)
	)
	return ret
}

func i32ToIP(a uint32) net.IP {
	return net.IPv4(byte(a>>24), byte(a>>16), byte(a>>8), byte(a))
}

func queryhttp(ip string) {
	httpx := NewHttpx(*urlFlag)
	httpx.SetTargetIP(ip)
	response, err := httpx.Send()
	if err != nil {
		fmt.Printf("fail : %s\n",  ip)
	}else{
		bodyString := string(response)
		if strings.Contains(bodyString, *wordFlag) {
			if *outputFlag != "" {
				if err := WriteFile(*outputFlag, 0600, fmt.Sprintf("%s\n", ip)); err != nil {
					fmt.Printf("write file error : %s", err)
				}
			}
			fmt.Printf("ok : %s\n", ip)		
		}

	}
}

func worker(linkChan chan string, wg *sync.WaitGroup) {
	// Decreasing internal counter for wait-group as soon as goroutine finishes
	defer wg.Done()

	for ip := range linkChan {
		queryhttp(ip)
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
		if strings.Index(*ipFlag, "/") != -1 {
			if _, _, err := net.ParseCIDR(*ipFlag); err == nil {
				cidr := strings.Split(*ipFlag, "/")
				addr32 := ipStringToI32(cidr[0])
				mask32, _ := strconv.ParseUint(cidr[1], 10, 8)
				ip_start := addr32 & (0xFFFFFFFF << (32 - mask32))
				ip_end := addr32 | ^(0xFFFFFFFF << (32 - mask32))
				fmt.Printf("ip_start :%d(%s) -> ip_end %d(%s)\n", ip_start, i32ToIP(ip_start).String(), ip_end, i32ToIP(ip_end).String())
				for ipint32 := ip_start; ipint32 <= ip_end; ipint32++ {
					lCh <- i32ToIP(ipint32).String()
				}
			} else {
				fmt.Fprintf(os.Stderr, "%s: Error %s\n", *ipFlag, err)
			}
		}

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
		for i, ip := range iplist {
			fmt.Printf("Start line %d/%d\n", i, len(iplist))
			ip = strings.Trim(ip, "\r\n ")
			if strings.Index(ip, "/") != -1 {
				if _, _, err := net.ParseCIDR(ip); err == nil {
					cidr := strings.Split(ip, "/")
					addr32 := ipStringToI32(cidr[0])
					mask32, _ := strconv.ParseUint(cidr[1], 10, 8)
					ip_start := addr32 & (0xFFFFFFFF << (32 - mask32))
					ip_end := addr32 | ^(0xFFFFFFFF << (32 - mask32))
					for ipint32 := ip_start; ipint32 <= ip_end; ipint32++ {
						lCh <- i32ToIP(ipint32).String()
					}
				}
			} else {
				lCh <- ip
			}
		}
	}

	// Closing channel (waiting in goroutines won't continue any more)
	close(lCh)

	// Waiting for all goroutines to finish (otherwise they die as main routine dies)
	wg.Wait()
}
