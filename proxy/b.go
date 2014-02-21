package main

import (
	"fmt"
	"net"
)

const (
	maxPkgSize = (4 * 1024 * 1024)
)

func start(netConn net.Conn, id int) {
	connInfo := new(gameConn)
	connInfo.ClientConn = netConn
	connInfo.Channel = make(chan int)

	connMap[id] = connInfo

	var err error
	connInfo.PkgLen, err = netConn.Read(connInfo.PkgBuf[:])
	if err != nil {
		fmt.Println(err)
		return
	} else if connInfo.PkgLen >= maxPkgSize {
		fmt.Printf("too long:%i\n", connInfo.PkgLen)
	} else {
		//		fmt.Printf("length:%i\n", length)
	}

	gbChannel <- id

	for {
		<-connInfo.Channel

		connInfo.PkgLen, err = netConn.Read(connInfo.PkgBuf[:])
		if err != nil {
			fmt.Println(err)
			return
		} else if connInfo.PkgLen >= maxPkgSize {
			fmt.Printf("too long:%i\n", connInfo.PkgLen)
		} else {
			//			fmt.Printf("length:%i\n", length)
		}

		gbChannel <- id
	}
}

func redisProcess(serverConn net.Conn) {
	var pGameConn *gameConn
	var id int
	for {
		id = <-gbChannel

		pGameConn = connMap[id]
		if pGameConn == nil {
			continue
		}

		length, err := serverConn.Write(pGameConn.PkgBuf[:pGameConn.PkgLen])
		if err != nil {
			fmt.Println(err)
			continue
		}

		length, err = serverConn.Read(pGameConn.PkgBuf[:])
		if err != nil {
			fmt.Println(err)
			continue
		} else if length >= maxPkgSize {
			fmt.Printf("too long:%i\n", length)
		} else {
			//		fmt.Printf("length:%i\n", length)
		}

		length, err = pGameConn.ClientConn.Write(pGameConn.PkgBuf[:length])
		if err != nil {
			fmt.Println(err)
			continue
		}

		pGameConn.Channel <- 1
	}

}

