package main

import (
    "bytes"
    "encoding/binary"
    . "fmt"
    "strconv"
    "strings"
    "syscall"
    "unsafe"
)

type IpHeader struct {
    Version  int         // protocol version
    Len      int         // header length
    TOS      int         // type-of-service
    TotalLen int         // packet total length
    ID       int         // identification
    //Flags    int // flags
    FragOff  int         // fragment offset
    TTL      int         // time-to-live
    Protocol int         // next protocol
    Checksum int         // checksum
    Src      net.IP      // source address
    Dst      net.IP      // destination address
    //Options  []byte      // options, extension headers
}

type UdpHdr struct {
	Source	uint16
	Dest	uint16
	Len	uint16
	Check 	uint16
}

type PsdHeader struct {
	SrcAddr	uint32
	DstAddr	uint32
	Filler	uint8
	Protocol	uint8
	Len	uint16
}
type DNSHeader struct {
	ID uint16
	Flag uint16
	QuestionCount uint16
	AnswerRRs uint16 //RRs is Resource Records
	AuthorityRRs uint16
	AdditionalRRs uint16
}

type Query struct
{
	Qtype uint16
	Qclass uint16
}

type TCPHeader struct {
    SrcPort   uint16
    DstPort   uint16
    SeqNum    uint32
    AckNum    uint32
    Offset    uint8
    Flag      uint8
    Window    uint16
    Checksum  uint16
    UrgentPtr uint16
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

func htons(port uint16) uint16 {
    var (
        high uint16 = port >> 8
        ret  uint16 = port<<8 + high
    )
    return ret
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

func main() {
    var (
        msg       string
        psdheader PsdHeader
        ipheader  IpHeader
        dnsheader DNSHeader
        query     Query
    )


	dnsheader.ID = 123
	dnsheader.Flag = htons(0x0100)
	dnsheader.QuestionCount = htons(1)
	dnsheader.AnswerRRs = 0
	dnsheader.AuthorityRRs = 0
	dnsheader.AdditionalRRs = 0

	query->Qtype = htons(0x00ff);
	query->Qclass = htons(0x1);
	
    /*填充UDP伪首部*/
    psdheader.SrcAddr = inet_addr("192.168.3.200")
    psdheader.DstAddr = inet_addr("123.125.114.144")
    psdheader.Filler = 0
    psdheader.Protocol = syscall.IPPROTO_UDP
    psdheader.Len = uint16(unsafe.Sizeof(UdpHdr{})) + uint16(len(msg))
    //htons(sizeof(udph) + sizeof(dns_hdr) + (strlen(dns_name)+1) + sizeof(query));

    /*填充IP首部*/
    ipheader.Version = 4
	ipheader.Len = 5         // header length
	ipheader.TOS = 0
	ipheader.TotalLen =	0
	ipheader.ID = 12345
	//ipheader.Flags:	DontFragment,
	ipheader.FragOff = 0
	ipheader.TTL = 64
	ipheader.Protocol = syscall.IPPROTO_UDP
	ipheader.Checksum =	0
	ipheader.Src = net.ParseIP("127.0.0.2")
	ipheader.Dst = net.ParseIP("127.0.0.1")
	//ipheader.Options = []byte{}

    /*buffer用来写入两种首部来求得校验和*/
    var (
        buffer bytes.Buffer
    )
    binary.Write(&buffer, binary.BigEndian, psdheader)
    binary.Write(&buffer, binary.BigEndian, ipheader)
    tcpheader.Checksum = CheckSum(buffer.Bytes())

    /*接下来清空buffer，填充实际要发送的部分*/
    buffer.Reset()
    binary.Write(&buffer, binary.BigEndian, tcpheader)
    binary.Write(&buffer, binary.BigEndian, msg)


	
    /*下面的操作都是raw socket操作，大家都看得懂*/
    var (
        sockfd int
        addr   syscall.SockaddrInet4
        err    error
    )
    if sockfd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_UDP); err != nil {
        Println("Socket() error: ", err.Error())
        return
    }
    defer syscall.Shutdown(sockfd, syscall.SHUT_RDWR)
    addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3] = 127, 0, 0, 0
    addr.Port = 53
    if err = syscall.Sendto(sockfd, buffer.Bytes(), 0, &addr); err != nil {
        Println("Sendto() error: ", err.Error())
        return
    }
    Println("Send success!")
}

