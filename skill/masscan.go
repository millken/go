package main

import (
	"encoding/xml"
	"flag"
	"log"
	"os"
	"strings"
	"fmt"
)

type Masscan struct {
	XMLName  xml.Name `xml:"nmaprun"`
	Scaninfo Scaninfo `xml:"scaninfo"`
	Host     []Host   `xml:"host"`
	Runstats Runstats `xml:"runstats"`
}

type Scaninfo struct {
	Type     string `xml:"type,attr"`
	Protocol string `xml:"protocol,attr"`
}

type Host struct {
	Address Address `xml:"address"`
	Ports   Ports   `xml:"ports"`
}

type Address struct {
	Addr     string `xml:"addr,attr"`
	Addrtype string `xml:"addrtype,attr"`
}

type Ports struct {
	Port    Port    `xml:"port"`
}

type Port struct {
	Protocol string `xml:"protocol,attr"`
	Portid   string `xml:"portid,attr"`
	State   State   `xml:"state"`
	Service Service `xml:"service"`	
}

type State struct {
	State  string `xml:"state,attr"`
	Reason string `xml:"reason,attr"`
	Reason_ttl int `xml:"reason_ttl,attr"`
}

type Service struct {
	Name   string `xml:"name,attr"`
	Banner string `xml:"banner,attr"`
}

type Runstats struct {
	Finished Finished `xml:"finished"`
	Hosts    Hosts    `xml:"hosts"`
}

type Finished struct {
	Time    string `xml:"time,attr"`
	Timestr string `xml:"timestr,attr"`
	Elapsed int    `xml:"elapsed,attr"`
}

type Hosts struct {
	Up    int `xml:"up,attr"`
	Down  int `xml:"down,attr"`
	Total int `xml:"total,attr"`
}

type Result struct {
	Servername, Title string
}

var (
	fFlag = flag.String("f", "", "ready masscan xml result file to parse")
	sFlag = flag.String("s", "", "match server name")
	tFlag = flag.String("t", "", "match title")
	oFlag = flag.String("o", "", "output file")
	total = 0
)

var masscan *Masscan
var resscan = make(map[string]*Result, 0)

func init() {
	flag.Parse()
	if *fFlag == "" {
		*fFlag = "config.xml"
	}
}

//http://shawnps.net/code/go/the-case-of-encoding-xml/
func ParseMasscan(filename string) (masscan Masscan) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == "nmaprun" {
				decoder.DecodeElement(&masscan, &startElement)

			}
		}
	}
	return
}

func getHeader(header *string, name string) (value string) {
	headers := strings.Split(*header, "\\x0d\\x0a")
	value = ""
	for _, hd := range headers {
		hds := strings.SplitN(hd, ": ", 2)
		if hds[0] == name {
			value = hds[1]
			return
		}
	}
	return
}

func setResscan(ip, name, value string) {
	if _, ok := resscan[ip]; !ok {
		resscan[ip] = &Result{"",""}
	}
	switch name {
		case "servername" :
			resscan[ip].Servername = value
		case "title" :
			resscan[ip].Title = value
	}
}

func main() {
	fmt.Fprintln(os.Stdout, "loading file...");
	masscan := ParseMasscan(*fFlag)
	var host Host
	var service Service
	for total, host = range masscan.Host {
		var ip, title, header, servername string
		ip = host.Address.Addr
		service = host.Ports.Port.Service
		if service.Name == "http" {
			header = service.Banner
			servername = getHeader(&header, "Server")
			setResscan(ip, "servername", servername)
		}
		if service.Name == "title" {
			title = service.Banner
			setResscan(ip, "title", title)
		}
		//fmt.Printf("%s (title)=%s, (server)=%s\n", ip, title, servername)
	}
	fmt.Printf("found %d record\n", total + 1)
	
	var rec = 1
	if *sFlag != "" { rec = rec | 2 }
	if *tFlag != "" { rec = rec | 4 }
		
	for ip, res := range resscan {
		recs := 1
		if *sFlag != "" && strings.Index(res.Servername, *sFlag) != -1 {recs = recs | 2}
		if *tFlag != "" && strings.Index(res.Title, *tFlag) != -1 {recs = recs | 4}
		if rec == recs {
			fmt.Printf("%s [%s] [%s]\n", ip, res.Servername, res.Title)
		}
	}
}
