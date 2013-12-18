package dns

////////////////////////////////////////////////////////////////////////////////
// Method functions
////////////////////////////////////////////////////////////////////////////////

func (dns *Connection) NewSimpleQuery(rtype RecordType, domain string) *Message {
	return &Message{
		Header: &Header{
			ID:            dns.cur_id,
			Query:         true,
			OpCode:        OPCODE_QUERY,
			Recursion:     true,
			QuestionCount: 1,
		},
		Questions: []*Question{
			NewQuestion(domain, rtype, CLASS_IN),
		},
	}
}
