package main

import "os"
import "fmt"
import "strconv"
import . "../../godns"

func main() {
	if len(os.Args) < 4 {
		println("Invalid args")
		println(os.Args[0], " [dns server] [domain] [type]")
		println("   A RECORD = 1")
		println(" TXT RECORD = 16")
		os.Exit(1)
	}

	conn, err := NewConnection(os.Args[1], 53)
	if err != nil {
		println("Error connecting:", err.Error())
		os.Exit(1)
	}
	rectype, err := strconv.Atoi(os.Args[3])
	if err != nil {
		println("Invalid record type:", err)
		os.Exit(1)
	}

	packet := conn.NewSimpleQuery(RecordType(rectype), os.Args[2])
	println(packet.String())
	resp, err := conn.Send(packet)
	if err != nil {
		println("Error sending:", err.Error())
		os.Exit(1)
	}
	fmt.Println(resp.String())
}
