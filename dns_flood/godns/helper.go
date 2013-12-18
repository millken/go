package dns

import "bytes"
import "encoding/binary"

func readDnsString(buffer *bytes.Buffer, buf []byte) (str string) {
	size := int((buffer.Next(1))[0])
	// Message pointer
	if size == 0xC0 {
		pointerb := buffer.Next(1)
		if len(pointerb) < 1 {
			return "[invalid pointer]"
		}
		pointer := pointerb[0]
		return readDnsString(bytes.NewBuffer(buf[pointer:]), buf)
	}
	// String
	for size != 0 {
		str += string(buffer.Next(size))
		size = int((buffer.Next(1))[0])
		if size != 0 {
			str += "."
		}
	}
	return
}

func write16(buf *bytes.Buffer, value uint16) {
	err := binary.Write(buf, binary.BigEndian, value)
	if err != nil {
		panic("write16 failed")
	}
}

func write8(buf *bytes.Buffer, value uint8) {
	err := binary.Write(buf, binary.BigEndian, value)
	if err != nil {
		panic("write16 failed")
	}
}
