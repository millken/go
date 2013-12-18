package dns

import "fmt"

////////////////////////////////////////////////////////////////////////////////
// Public functions
////////////////////////////////////////////////////////////////////////////////

func NewHeader() *Header {
	return &Header{}
}

func NewQueryHeader(id int, recursion bool) *Header {
	return &Header{}
}

////////////////////////////////////////////////////////////////////////////////
// Method functions
////////////////////////////////////////////////////////////////////////////////

func (opcode OpCode) String() string {
	switch opcode {
	case OPCODE_QUERY:
		return "query"
	case OPCODE_IQUERY:
		return "iquery"
	case OPCODE_STATUS:
		return "status"
	case OPCODE_UNASSIGNED:
		return "unassigned"
	case OPCODE_NOTIFY:
		return "notify"
	case OPCODE_UPDATE:
		return "update"
	}
	return "Unknown opcode"
}

func (rcode ResponseCode) String() string {
	switch rcode {
	case RCODE_NOERROR:
		return "ok"
	case RCODE_FORMAT_ERROR:
		return "format error"
	case RCODE_SERVER_FAILURE:
		return "server failure"
	case RCODE_NON_EXISTANT_DOMAIN:
		return "non-existant domain"
	case RCODE_NOT_IMPLEMENTED:
		return "not implemented"
	case RCODE_QUERY_REFUSED:
		return "query refused"
	}
	return "Unknown respone code"
}

func (h *Header) String() (str string) {
	str += fmt.Sprintf("                ID: %d\n", h.ID)
	str += fmt.Sprintf("             Query: %v\n", h.Query)
	str += fmt.Sprintf("         Operation: %s\n", h.OpCode.String())
	str += fmt.Sprintf("     Authoritative: %v\n", h.Authoritative)
	str += fmt.Sprintf("         Truncated: %v\n", h.Truncated)
	str += fmt.Sprintf("         Recursion: %v\n", h.Recursion)
	str += fmt.Sprintf("RecursionSupported: %v\n", h.RecursionSupported)
	str += fmt.Sprintf("          Reserved: %v\n", h.Reserved)
	str += fmt.Sprintf("          Response: %s\n", h.ResponseCode.String())
	str += fmt.Sprintf("     QuestionCount: %v\n", h.QuestionCount)
	str += fmt.Sprintf("       AnswerCount: %v\n", h.AnswerCount)
	str += fmt.Sprintf("   NameserverCount: %v\n", h.NameserverCount)
	str += fmt.Sprintf("   AdditionalCount: %v\n", h.AdditionalCount)
	return
}
