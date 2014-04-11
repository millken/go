/*

*/

package main

import (
	"./godns"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var (
	nConnectFlag = flag.Int("concurrent", 100, "Number of concurrent connections")
	dataFileFlag = flag.String("data", "", "File containing message to send to responsive hosts ('%s' will be replaced with host IP)")
	domainFlag   = flag.String("domain", "www.google.com", "checkdomain, default is : www.google.com")
	ipFlag       = flag.String("ip", "", "ip address(es) cidr to scan. example 8.8.8.8/24")
	outputFlag   = flag.String("output", "", "output file")
)

var messageData = make([]byte, 0) // data read from file specified with dataFile flag

// Before running main, parse flags and load message data, if applicable
func init() {
	flag.Parse()
	// Increase file descriptor limit
	rlimit := syscall.Rlimit{Max: uint64(*nConnectFlag + 50000), Cur: uint64(*nConnectFlag + 50000)}
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting rlimit: %s", err)
	}
}

type Work struct {
	ip string
}

type resultStruct struct {
	addr string   // address of remote host
	data []string // data returned from the host, if successful
	err  error    // error, if any
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

func worker(id int, queue chan *Work, resultChan chan resultStruct) {
	for {
		wp := <-queue
		if wp == nil {
			break
		}
		//fmt.Printf("worker #%d: item %v\n", id, *wp)
		handleWork(wp, resultChan)
	}
}
func handleWork(w *Work, resultChan chan resultStruct) {
	//fmt.Printf("handleWork  %v\n", w)
	options := &godns.LookupOptions{
		DNSServers: []string{w.ip},
		Net:        "udp",
		CacheTTL:   godns.DNS_NOCACHE,
		OnlyIPv4:   true,
		//DialTimeout: DialTimeout("udp")
	}
	addrs, err := godns.LookupHost(*domainFlag, options)
	if err != nil {
		resultChan <- resultStruct{w.ip, nil, err}
	} else {
		ret := []string{}
		for _, ip := range addrs {
			ret = append(ret, ip)
		}
		resultChan <- resultStruct{w.ip, ret, nil}
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

func output(resultChan chan resultStruct, doneChan chan int) {

	ok, error := 0, 0
	for result := range resultChan {
		if result.err == nil {
			fmt.Printf("%s: %v\n", result.addr, result.data)
			if *outputFlag != "" {
				if err := WriteFile(*outputFlag, 0600, fmt.Sprintf("%s\n", result.addr));err != nil {
					fmt.Printf("write file error : %s", err)
				}
			}
			ok++
		} else {
			//fmt.Fprintf(os.Stderr, "%s: Error %s\n", result.addr, result.err)
			error++
		}
	}
	fmt.Fprintf(os.Stderr, "Complete %d (success=%d, failure=%d)\n", ok + error, ok, error)
	doneChan <- 1
}

func main() {
	workChan := make(chan *Work, *nConnectFlag)
	resultChan := make(chan resultStruct, *nConnectFlag) // grabbers send results to output
	doneChan := make(chan int, *nConnectFlag)            // let grabbers signal completion	
	

	for i := 0; i < *nConnectFlag; i++ {
		go worker(i, workChan, resultChan)
	}
	
	go output(resultChan, doneChan)
	
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
					workChan <- &Work{ i32ToIP(ipint32).String() }
				}
			} else {
				fmt.Fprintf(os.Stderr, "%s: Error %s\n", *ipFlag, err)
			}
		}

	}
	for n := 0; n < *nConnectFlag; n++ {
        workChan <- nil
    }

	close(resultChan)

}
