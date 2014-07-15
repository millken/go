/*
http://blog.csdn.net/gophers/article/details/22942457
http://blog.csdn.net/gophers/article/details/20393601
https://github.com/grahamking/latency/blob/master/latency.go
http://www.51testing.com/html/66/138366-216709.html IP包头结构详解
*/
package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"syscall"
	"strconv"
	"strings"
	"net"
	"unsafe"
)

type PsdHeader struct {
	SrcAddr	uint32
	DstAddr	uint32
	Filler	uint8
	Protocol	uint8
	Len	uint16
}

type IpHeader struct {
    //Version  uint8         // protocol version
    //Len      uint8         // header length
    VerLen	 uint8
    TOS      uint8         // type-of-service
    TotalLen uint16         // packet total length
    ID       uint16         // identification
    //Flags    uint8 // flags
    FragOff  uint16         // fragment offset
    TTL      uint8         // time-to-live
    Protocol uint8         // next protocol
    Checksum uint16         // checksum
    Src      uint32      // source address
    Dst      uint32      // destination address
    //Options  []byte      // options, extension headers
    //Padding  []byte
}

type UdpHdr struct {
	Source	uint16
	Dest	uint16
	Len	uint16
	Check 	uint16
}

type DNSHeader struct {
	ID            uint16
	Flag          uint16
	QuestionCount uint16
	AnswerRRs     uint16 //RRs is Resource Records
	AuthorityRRs  uint16
	AdditionalRRs uint16
}

func (header *DNSHeader) SetFlag(QR uint16, OperationCode uint16, AuthoritativeAnswer uint16, Truncation uint16, RecursionDesired uint16, RecursionAvailable uint16, ResponseCode uint16) {
	header.Flag = QR<<15 + OperationCode<<11 + AuthoritativeAnswer<<10 + Truncation<<9 + RecursionDesired<<8 + RecursionAvailable<<7 + ResponseCode
}

type DNSQuery struct {
	QuestionType  uint16
	QuestionClass uint16
}

func inet_addr(ipaddr string) uint32 {
    var (
        segments []string = strings.Split(ipaddr, ".")
        ip       [4]uint64
        ret      uint64
    )
    for i := 0; i < 4; i++ {
        ip[i], _ = strconv.ParseUint(segments[i], 10, 64)
    }
    ret = ip[3]<<24 + ip[2]<<16 + ip[1]<<8 + ip[0]
    return uint32(ret)
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

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

func ParseDomainName(domain string) []byte {
	//要将域名解析成相应的格式，例如：
	//"www.google.com"会被解析成"0x03www0x06google0x03com0x00"
	//就是长度+内容，长度+内容……最后以0x00结尾
	var (
		buffer   bytes.Buffer
		segments []string = strings.Split(domain, ".")
	)
	for _, seg := range segments {
		binary.Write(&buffer, binary.BigEndian, byte(len(seg)))
		binary.Write(&buffer, binary.BigEndian, []byte(seg))
	}
	binary.Write(&buffer, binary.BigEndian, byte(0x00))

	return buffer.Bytes()
}

func main() {
	var (
	    psdheader    PsdHeader
		ipheader     IpHeader
		udphdr		 UdpHdr
		dns_header   DNSHeader
		dns_question DNSQuery
	)

	domain := "www.baidu.com"
	format_domain := ParseDomainName(domain)

	//udphdr_len := 8

    psdheader.SrcAddr = ipStringToI32("8.8.8.8")
    psdheader.DstAddr = ipStringToI32("127.0.0.1")
    psdheader.Filler = 0
    psdheader.Protocol = syscall.IPPROTO_UDP
    psdheader.Len = uint16(unsafe.Sizeof(UdpHdr{})) + uint16(unsafe.Sizeof(DNSHeader{})) + uint16(len(format_domain)) + uint16(unsafe.Sizeof(DNSQuery{}))


    ipheader.VerLen = 0x45
	ipheader.TOS = 0x00
	ipheader.TotalLen = uint16(unsafe.Sizeof(IpHeader{})) + uint16(unsafe.Sizeof(UdpHdr{})) + uint16(unsafe.Sizeof(DNSHeader{})) + uint16(len(format_domain)) + uint16(unsafe.Sizeof(DNSQuery{}))

	ipheader.ID = 12345
	ipheader.FragOff = 0
	ipheader.TTL = 64
	ipheader.Protocol = syscall.IPPROTO_UDP
	ipheader.Checksum =	0
	ipheader.Src = ipStringToI32("8.8.8.8")
	ipheader.Dst = ipStringToI32("127.0.0.1")
	

	udphdr.Source = 0x74b8
	udphdr.Dest = 53
	udphdr.Len = uint16(8) + uint16(unsafe.Sizeof(DNSHeader{})) + uint16(len(format_domain)) + uint16(unsafe.Sizeof(DNSQuery{}))
	udphdr.Check = 0
				
	//填充dns首部
	dns_header.ID = 0x3f58
	dns_header.SetFlag(0, 0, 0, 0, 1, 0, 0)
	dns_header.QuestionCount = 1
	dns_header.AnswerRRs = 0
	dns_header.AuthorityRRs = 0
	dns_header.AdditionalRRs = 0

	//填充dns查询首部
	dns_question.QuestionType = 1  //IPv4
	dns_question.QuestionClass = 1  

	var (
		addr   syscall.SockaddrInet4
		err  error
		sockfd int
		buffer bytes.Buffer
	)

	binary.Write(&buffer, binary.BigEndian, psdheader)
	binary.Write(&buffer, binary.BigEndian, ipheader)

    ipheader.Checksum = CheckSum(buffer.Bytes())
    
    /*接下来清空buffer，填充实际要发送的部分*/
    buffer.Reset()   
    //binary.Write(&buffer, binary.BigEndian, psdheader) 	
	//binary.Write(&buffer, binary.BigEndian, ipheader)
	binary.Write(&buffer, binary.BigEndian, udphdr)
	//buffer中是我们要发送的数据，里面的内容是DNS首部+查询内容+DNS查询首部
	binary.Write(&buffer, binary.BigEndian, dns_header)
	binary.Write(&buffer, binary.BigEndian, format_domain)
	binary.Write(&buffer, binary.BigEndian, dns_question)
	fmt.Printf("%s: %s\nlength=%d\n", buffer.Bytes(),
					hex.EncodeToString(buffer.Bytes()), len(buffer.Bytes()))

	
    if sockfd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_UDP); err != nil {
        fmt.Println("Socket() error: ", err.Error())
        return
    }
    defer syscall.Shutdown(sockfd, syscall.SHUT_RDWR)
    addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3] = 127, 0, 0, 1
    addr.Port = 53
    if err = syscall.Sendto(sockfd, buffer.Bytes(), 0, &addr); err != nil {
        fmt.Println("Sendto() error: ", err.Error())
        return
    }	
    fmt.Println("Send success!")
}

