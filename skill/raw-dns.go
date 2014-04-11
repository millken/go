/*
http://blog.csdn.net/gophers/article/details/20393601
https://github.com/grahamking/latency/blob/master/latency.go
http://www.51testing.com/html/66/138366-216709.html IP包头结构详解
*/
package main
import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)
const (
        Version      = 4  // protocol version
        HeaderLen    = 20 // header length without extension headers
        maxHeaderLen = 60 // sensible default, revisit if later RFCs define new usage of version and header length fields
)
const (
        posTOS      = 1  // type-of-service
        posTotalLen = 2  // packet total length
        posID       = 4  // identification
        posFragOff  = 6  // fragment offset
        posTTL      = 8  // time-to-live
        posProtocol = 9  // next protocol
        posChecksum = 10 // checksum
        posSrc      = 12 // source address
        posDst      = 16 // destination address
)

const supportsNewIPInput = runtime.GOOS == "linux" || runtime.GOOS == "openbsd"

var (
	domainFlag   = flag.String("domain", "www.google.com", "checkdomain, default is : www.google.com")
)
type HeaderFlags int

const (
        MoreFragments HeaderFlags = 1 << iota // more fragments flag
        DontFragment                          // don't fragment flag
)
type IpHeader struct {
    Version  int         // protocol version
    Len      int         // header length
    TOS      int         // type-of-service
    TotalLen int         // packet total length
    ID       int         // identification
    Flags    int // flags
    FragOff  int         // fragment offset
    TTL      int         // time-to-live
    Protocol int         // next protocol
    Checksum int         // checksum
    Src      net.IP      // source address
    Dst      net.IP      // destination address
    Options  []byte      // options, extension headers
}

type UdpHdr struct {
	Source	int
	Dest	int
	Len	int
	Check 	int
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

func (header *DNSHeader) SetFlag(QR uint16, OperationCode uint16, AuthoritativeAnswer uint16, Truncation uint16, RecursionDesired uint16, RecursionAvailable uint16, ResponseCode uint16) {
	header.Flag = QR<<15 + OperationCode<<11 + AuthoritativeAnswer<<10 + Truncation<<9 + RecursionDesired<<8 + RecursionAvailable<<7 + ResponseCode
}

type DNSQuery struct {
	QuestionType uint16
	QuestionClass uint16
}
func (h *IpHeader) Marshal() ([]byte, error) {
        if h == nil {
                return nil, syscall.EINVAL
        }
        if h.Len < HeaderLen {
                return nil, errHeaderTooShort
        }
        hdrlen := HeaderLen + len(h.Options)
        b := make([]byte, hdrlen)
        b[0] = byte(Version<<4 | (hdrlen >> 2 & 0x0f))
        b[posTOS] = byte(h.TOS)
        flagsAndFragOff := (h.FragOff & 0x1fff) | int(h.Flags<<13)
        if supportsNewIPInput {
                b[posTotalLen], b[posTotalLen+1] = byte(h.TotalLen>>8), byte(h.TotalLen)
                b[posFragOff], b[posFragOff+1] = byte(flagsAndFragOff>>8), byte(flagsAndFragOff)
        } else {
                *(*uint16)(unsafe.Pointer(&b[posTotalLen : posTotalLen+1][0])) = uint16(h.TotalLen)
                *(*uint16)(unsafe.Pointer(&b[posFragOff : posFragOff+1][0])) = uint16(flagsAndFragOff)
        }
        b[posID], b[posID+1] = byte(h.ID>>8), byte(h.ID)
        b[posTTL] = byte(h.TTL)
        b[posProtocol] = byte(h.Protocol)
        b[posChecksum], b[posChecksum+1] = byte(h.Checksum>>8), byte(h.Checksum)
        if ip := h.Src.To4(); ip != nil {
                copy(b[posSrc:posSrc+net.IPv4len], ip[:net.IPv4len])
        }
        if ip := h.Dst.To4(); ip != nil {
                copy(b[posDst:posDst+net.IPv4len], ip[:net.IPv4len])
        } else {
                return nil, errMissingAddress
        }
        if len(h.Options) > 0 {
                copy(b[HeaderLen:], h.Options)
        }
        return b, nil
}
/*
// Pseudoheader struct
typedef struct
{
    u_int32_t saddr;
    u_int32_t daddr;
    u_int8_t filler;
    u_int8_t protocol;
    u_int16_t len;
}ps_hdr;

// DNS header struct
typedef struct
{
	unsigned short id; 		// ID
	unsigned short flags;	// DNS Flags
	unsigned short qcount;	// Question Count
	unsigned short ans;		// Answer Count
	unsigned short auth;	// Authority RR
	unsigned short add;		// Additional RR
}dns_hdr;

// Question types
typedef struct
{
	unsigned short qtype;
	unsigned short qclass;
}query;
*/

func init() {
	flag.Parse()
	// Increase file descriptor limit
	rlimit := syscall.Rlimit{Max: uint64(500000), Cur: uint64(500000)}
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting rlimit: %s", err)
	}
}

func main() {
	ipheader := IpHeader{
		Version:	4,
		Len      0,         // header length
		TOS:	0,
		TotalLen:	0,
		ID:	12345,
		Flags:	DontFragment,
		FragOff:	0,
		TTL:	64,
		Protocol:	syscall.IPPROTO_UDP,
		Checksum:	0,
		Src:      net.ParseIP("127.0.0.2"),
		Dst:	net.ParseIP("127.0.0.1"),
		Options  []byte{}
	}
	psdheader := PsdHeader{
		SrcAddr:	inet_addr("127.0.0.2"),
		DstAddr:	inet_addr("127.0.0.1"),
		Filler:	0,
		Protocol:	syscall.IPPROTO_UDP,
		Len:	htons(sizeof(udph) + sizeof(dns_hdr) + (strlen(dns_name)+1) + sizeof(query)),
	}
}
