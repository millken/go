/*

*/

package main

import (
	"./godns"
	"bufio"
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
)

var messageData = make([]byte, 0) // data read from file specified with dataFile flag

// Before running main, parse flags and load message data, if applicable
func init() {
	flag.Parse()
	// Increase file descriptor limit
	rlimit := syscall.Rlimit{Max: uint64(*nConnectFlag + 4), Cur: uint64(*nConnectFlag + 4)}
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

// Read addresses from addrChan and grab banners from these hosts.
// Sends resultStructs to resultChan.  Writes to doneChan when complete.
func grabber(addrChan chan string, resultChan chan resultStruct, doneChan chan int) {
	for{
	select {
		case addr := <-addrChan :
		options := &godns.LookupOptions{
			DNSServers: []string{addr},
			Net:        "udp",
			CacheTTL:   godns.DNS_NOCACHE,
			OnlyIPv4:   true,
			//DialTimeout: DialTimeout("udp")
		}
		addrs, err := godns.LookupHost(*domainFlag, options)
		if err != nil {
			resultChan <- resultStruct{addr, nil, err}
		} else {
			ret := []string{}
			for _, ip := range addrs {
				ret = append(ret, ip)
			}
			resultChan <- resultStruct{addr, ret, nil}
		}
	}
	}
	doneChan <- 1
}

// Read resultStructs from resultChan, print output, and maintain
// status counters.  Writes to doneChan when complete.
func output(resultChan chan resultStruct, doneChan chan int) {
	ok, error := 0, 0
	for result := range resultChan {
		if result.err == nil {
			fmt.Printf("%s: %v\n", result.addr,
				result.data)
			ok++
		} else {
			//fmt.Fprintf(os.Stderr, "%s: Error %s\n", result.addr, result.err)
			error++
		}
	}
	fmt.Fprintf(os.Stderr, "Complete %d (success=%d, failure=%d)\n", ok + error, ok, error)
	doneChan <- 1
}

//http://play.golang.org/p/TZbIBev4pU

func ipStringToI32(a string) uint32 {
	return ipToI32(net.ParseIP(a))
}
func ipToI32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func i32ToIP(a uint32) net.IP {
	return net.IPv4(byte(a>>24), byte(a>>16), byte(a>>8), byte(a))
}

func worker() {
}

func main() {
	addrChan := make(chan string, *nConnectFlag)     // pass addresses to grabbers
	resultChan := make(chan resultStruct, *nConnectFlag) // grabbers send results to output
	doneChan := make(chan int, *nConnectFlag)            // let grabbers signal completion

	// Start grabbers and output thread
	go output(resultChan, doneChan)
	for i := 0; i < *nConnectFlag; i++ {
		go grabber(addrChan, resultChan, doneChan)
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
					addrChan <-i32ToIP(ipint32).String()
				}
			} else {
				fmt.Fprintf(os.Stderr, "%s: Error %s\n", *ipFlag, err)
			}
		}

	} else {
		// Read addresses from stdin and pass to grabbers
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			addrChan <- scanner.Text()
		}
	}
	close(addrChan)

	// Wait for completion
	for i := 0; i < *nConnectFlag; i++ {
		<-doneChan
	}
	close(resultChan)
	<-doneChan
}
