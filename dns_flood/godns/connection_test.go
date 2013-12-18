package dns

import "testing"

const TEST_DNS_HOST = "localhost"
const TEST_DNS_PORT = 53

func testNewConnection(t *testing.T, host string, port int) (conn *Connection) {
	conn, err := NewConnection(host, port)
	if err != nil {
		t.Fatal(err)
	}
	if conn == nil {
		t.Fatal("NewConnection returned nil and no error!")
	}
	return
}

func TestNewConnection(t *testing.T) {
	testNewConnection(t, TEST_DNS_HOST, TEST_DNS_PORT)
}

func TestNewConnectionNeg(t *testing.T) {
	testNewConnection(t, "localhost", 9999)
}

func TestNewSimpleQuery(t *testing.T) {
	expected := []byte{0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 6, 109, 108, 111, 99, 97, 108, 4, 106, 111, 115, 104, 3, 99, 111, 109, 0, 0, 1, 0, 1}
	conn := testNewConnection(t, TEST_DNS_HOST, TEST_DNS_PORT)
	packet := conn.NewSimpleQuery(RECORD_TYPE_A, "mlocal.josh.com")
	if string(packet.Bytes()) != string(expected) {
		t.Error("Got:     ", packet.Bytes())
		t.Error("Expected:", expected)
		t.Fail()
	}
}

/*
func TestNewSimpleQueryIRL(t *testing.T) {
	dns := NewDns("8.8.8.8", 53)
	packet := dns.NewSimpleQuestion("mague.com")
	println(packet.String())
	resp, err := dns.Send(packet)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	println(resp)
	resp_packet := ParsePacket(resp)
	println(resp_packet.String())
}
*/
func TestNewTextQueryIRL(t *testing.T) {
	conn := testNewConnection(t, TEST_DNS_HOST, TEST_DNS_PORT)
	packet := conn.NewSimpleQuery(RECORD_TYPE_TXT, "www.fcsak.com")
	t.Log("\n", packet)
	resp, err := conn.Send(packet)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log("\n", resp)
}
