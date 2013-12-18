package dns

import "net"
import "fmt"

//import "os"

type Error error

//type Error os.Error

type Connection struct {
	cur_id uint16
	net.Conn
}

////////////////////////////////////////////////////////////////////////////////
// Public functions
////////////////////////////////////////////////////////////////////////////////

func NewConnection(server string, port int) (conn *Connection, err Error) {
	// Try to connect to the server
	udpconn, err := net.Dial("udp", fmt.Sprint(server, ":", port))
	if err != nil {
		return
	}

	// Connected, create an object
	conn = &Connection{
		cur_id: 1,
		Conn:   udpconn,
	}

	return
}

////////////////////////////////////////////////////////////////////////////////
// Method functions
////////////////////////////////////////////////////////////////////////////////

func (conn *Connection) Send(message *Message) (resp *Message, err Error) {
	// Write the request to the connection
	_, err = conn.Write(message.Bytes())
	if err != nil {
		return
	}

	// Read a response (Max message size 2000)
	buf := make([]byte, 2000)
	length, err := conn.Read(buf)
	if err != nil {
		return
	}

	// Trim the buffer and parse a message
	buf = buf[0:length]
	resp, err = ParseMessage(buf)

	return
}
