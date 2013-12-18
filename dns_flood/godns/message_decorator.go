package dns

////////////////////////////////////////////////////////////////////////////////
// Method functions
////////////////////////////////////////////////////////////////////////////////

func (message *Message) String() (str string) {
	if message.Query {
		str += "========== DNS Query ==========\n"
	} else {
		str += "========= DNS Response ========\n"
	}
	str += message.Header.String()
	if int(message.QuestionCount) > 0 {
		str += "======= Questions =======\n"
	}
	for i := 0; i < int(message.QuestionCount); i++ {
		str += message.Questions[i].String()
	}
	if int(message.AnswerCount) > 0 {
		str += "======== Answers ========\n"
	}
	for i := 0; i < int(message.AnswerCount); i++ {
		str += message.Answers[i].String()
	}
	if int(message.NameserverCount) > 0 {
		str += "====== Nameservers ======\n"
	}
	for i := 0; i < int(message.NameserverCount); i++ {
		str += message.Nameservers[i].String()
	}
	if int(message.AdditionalCount) > 0 {
		str += "====== Additionals ======\n"
	}
	for i := 0; i < int(message.AdditionalCount); i++ {
		str += message.Additionals[i].String()
	}
	return
}
