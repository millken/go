package dns

import "bytes"
import "encoding/binary"

type Answer struct {
	Name     string
	Type     RecordType
	Class    ClassType
	TTL      uint32
	DataSize uint16
	Data     []byte
}

////////////////////////////////////////////////////////////////////////////////
// Public functions
////////////////////////////////////////////////////////////////////////////////

func ParseAnswer(buffer *bytes.Buffer, buf []byte) (answer *Answer) {
	// Create a new answer object
	answer = NewAnswer()

	// Read the answer name (dns string or pointer)
	answer.Name = readDnsString(buffer, buf)

	// Parse 4 big endian words from the buffer
	binary.Read(buffer, binary.BigEndian, &answer.Type)
	binary.Read(buffer, binary.BigEndian, &answer.Class)
	binary.Read(buffer, binary.BigEndian, &answer.TTL)
	binary.Read(buffer, binary.BigEndian, &answer.DataSize)

	// The remaining data is Type specific
	answer.Data = buffer.Next(int(answer.DataSize))

	return
}

////////////////////////////////////////////////////////////////////////////////
// Method functions
////////////////////////////////////////////////////////////////////////////////

func (answer *Answer) Bytes() []byte {
	// TODO: Implement marshal
	return []byte{0}
}
