package dns

import "bytes"
import "strings"
import "encoding/binary"

////////////////////////////////////////////////////////////////////////////////
// Types
////////////////////////////////////////////////////////////////////////////////

type RecordType uint16

const (
	RECORD_TYPE_A RecordType = iota + 1
	RECORD_TYPE_NS
	RECORD_TYPE_MD
	RECORD_TYPE_MF
	RECORD_TYPE_CNAME
	RECORD_TYPE_SOA
	RECORD_TYPE_MB
	RECORD_TYPE_MG
	RECORD_TYPE_MR
	RECORD_TYPE_NULL
	RECORD_TYPE_WKS
	RECORD_TYPE_PTR
	RECORD_TYPE_HINFO
	RECORD_TYPE_MINFO
	RECORD_TYPE_MX
	RECORD_TYPE_TXT
)

type ClassType uint16

const (
	CLASS_IN ClassType = 1
)

type Question struct {
	Name  string
	Type  RecordType
	Class ClassType
}

////////////////////////////////////////////////////////////////////////////////
// Public functions
////////////////////////////////////////////////////////////////////////////////

func ParseQuestion(buffer *bytes.Buffer, buf []byte) (q *Question) {
	q = &Question{}
	q.Name = readDnsString(buffer, buf)
	binary.Read(buffer, binary.BigEndian, &q.Type)
	binary.Read(buffer, binary.BigEndian, &q.Class)
	return q
}

////////////////////////////////////////////////////////////////////////////////
// Method functions
////////////////////////////////////////////////////////////////////////////////

func (q *Question) Bytes() []byte {
	buf := new(bytes.Buffer)
	parts := strings.Split(q.Name, ".")
	for _, part := range parts {
		write8(buf, uint8(len(part)))
		buf.Write([]byte(part))
	}
	write8(buf, 0)
	write16(buf, uint16(q.Type))
	write16(buf, uint16(q.Class))
	return buf.Bytes()
}
