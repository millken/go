package dns

import "bytes"
import "encoding/binary"

////////////////////////////////////////////////////////////////////////////////
// Types
////////////////////////////////////////////////////////////////////////////////

type OpCode uint8

const (
	OPCODE_QUERY  OpCode = iota // RFC1035
	OPCODE_IQUERY               // RFC3425 (Obsolete)
	OPCODE_STATUS               // RFC1035
	OPCODE_UNASSIGNED
	OPCODE_NOTIFY // RFC1996
	OPCODE_UPDATE // RFC2136
)

type ResponseCode uint8

const (
	RCODE_NOERROR ResponseCode = iota
	RCODE_FORMAT_ERROR
	RCODE_SERVER_FAILURE
	RCODE_NON_EXISTANT_DOMAIN
	RCODE_NOT_IMPLEMENTED
	RCODE_QUERY_REFUSED
)

type Header struct {
	ID                 uint16       // WORD: Message ID
	Query              bool         // BIT 7: 0:Query, 1:Response
	OpCode             OpCode       // BIT 6-3: 1: Standard Query
	Authoritative      bool         // BIT 2: Authoratative answer 1:true 0:false
	Truncated          bool         // BIT 1: 1:message truncated, 0:normal
	Recursion          bool         // BIT 0: 1:request recursion, 0:no recursion
	RecursionSupported bool         // BIT 7: 1:recursion supported, 0:not
	Reserved           uint8        // BIT 6-4: Reserved 0
	ResponseCode       ResponseCode // BIT 3-0: Response code
	QuestionCount      uint16       // WORD: Number of queries
	AnswerCount        uint16       // WORD:
	NameserverCount    uint16       // WORD:
	AdditionalCount    uint16       // WORD:
}

////////////////////////////////////////////////////////////////////////////////
// Private functions
////////////////////////////////////////////////////////////////////////////////

func (h *Header) parse_byte3(byte3 byte) {
	if (byte3 & 0x80) == 0 {
		h.Query = true
	}
	h.OpCode = OpCode((byte3 >> 3) & 0xF)
	if (byte3>>2)&1 == 1 {
		h.Authoritative = true
	}
	if (byte3>>1)&1 == 1 {
		h.Truncated = true
	}
	if byte3&1 == 1 {
		h.Recursion = true
	}
}

func (h *Header) byte3() byte {
	var byte3 uint8
	if !h.Query {
		byte3 |= 0x80
	}
	byte3 |= (uint8(h.OpCode) & 0xF) << 3
	if h.Authoritative {
		byte3 |= 0x04
	}
	if h.Truncated {
		byte3 |= 0x02
	}
	if h.Recursion {
		byte3 |= 0x01
	}
	return byte3
}

func (h *Header) parse_byte4(byte4 byte) {
	if (byte4 & 0x80) != 0 {
		h.RecursionSupported = true
	}
	h.Reserved = (byte4 >> 4) & 0x7
	h.ResponseCode = ResponseCode(byte4 & 0x0F)
}

func (h *Header) byte4() byte {
	var byte4 uint8
	if h.RecursionSupported {
		byte4 |= 0x80
	}
	byte4 |= (h.Reserved & 0x7) << 4
	byte4 |= (uint8(h.ResponseCode) & 0x0F)
	return byte4
}

////////////////////////////////////////////////////////////////////////////////
// Public functions
////////////////////////////////////////////////////////////////////////////////

func ParseHeader(buffer *bytes.Buffer) *Header {
	h := &Header{}
	binary.Read(buffer, binary.BigEndian, &h.ID)
	var byte3 uint8
	binary.Read(buffer, binary.BigEndian, &byte3)
	h.parse_byte3(byte3)
	var byte4 uint8
	binary.Read(buffer, binary.BigEndian, &byte4)
	h.parse_byte4(byte4)
	binary.Read(buffer, binary.BigEndian, &h.QuestionCount)
	binary.Read(buffer, binary.BigEndian, &h.AnswerCount)
	binary.Read(buffer, binary.BigEndian, &h.NameserverCount)
	binary.Read(buffer, binary.BigEndian, &h.AdditionalCount)
	return h
}

////////////////////////////////////////////////////////////////////////////////
// Method functions
////////////////////////////////////////////////////////////////////////////////

func (h *Header) Bytes() []byte {
	buf := new(bytes.Buffer)
	write16(buf, h.ID)
	write8(buf, h.byte3())
	write8(buf, h.byte4())
	write16(buf, h.QuestionCount)
	write16(buf, h.AnswerCount)
	write16(buf, h.NameserverCount)
	write16(buf, h.AdditionalCount)
	return buf.Bytes()
}
