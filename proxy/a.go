package main

import (
	"fmt"
	"log"
	"net"
)

const (
	realServer = "127.0.0.1:9002"
	selfPort   = 9010
)

type gameConn struct {
	ClientConn net.Conn
	PkgBuf     [maxPkgSize]byte
	PkgLen     int
	Channel    chan int
}

var gbChannel chan int
var connMap map[int](*gameConn)

func main() {
	// Listen on TCP port 2000 on all interfaces.
	ip := net.ParseIP("0.0.0.0")
	addr := net.TCPAddr{ip, selfPort}
	l, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		log.Fatal(err)
	}

	serverConn, err := net.Dial("tcp", realServer)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("starting server...")

	gbChannel = make(chan int)
	connMap = make(map[int](*gameConn))
	index := 0

	go redisProcess(serverConn)

	for {
		// Wait for a connection.
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}

		go start(conn, index)

		index++
	}
}
