package dns

import "testing"

func compareResult(t *testing.T, got []byte, expected []byte) {
	if string(got) != string(expected) {
		t.Error("Got:", got)
		t.Error("Expected:", expected)
		t.Fail()
	}
}

func TestQuestionA(t *testing.T) {
	expected := []byte{3, 119, 119, 119, 12, 110, 111, 114, 116, 104, 101, 97, 115, 116, 101, 114, 110, 3, 101, 100, 117, 0, 0, 1, 0, 1}
	q := NewQuestion("www.northeastern.edu", RECORD_TYPE_A, CLASS_IN)
	compareResult(t, q.Bytes(), expected)
}

func TestQuestionTxt(t *testing.T) {
	expected := []byte{6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 16, 0, 1}
	q := NewQuestion("google.com", RECORD_TYPE_TXT, CLASS_IN)
	compareResult(t, q.Bytes(), expected)
}

func TestQuestionCname(t *testing.T) {
	expected := []byte{6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 5, 0, 1}
	q := NewQuestion("google.com", RECORD_TYPE_CNAME, CLASS_IN)
	compareResult(t, q.Bytes(), expected)
}

func TestQuestionMx(t *testing.T) {
	expected := []byte{4, 109, 97, 105, 108, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 15, 0, 1}
	q := NewQuestion("mail.google.com", RECORD_TYPE_MX, CLASS_IN)
	compareResult(t, q.Bytes(), expected)
}
