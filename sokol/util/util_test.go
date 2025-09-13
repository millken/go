package util_test

import (
	"slices"
	"testing"
	"unsafe"

	"github.com/millken/go/sokol/util"
)

func TestToCString(t *testing.T) {
	cases := []struct {
		Given  string
		Wanted []byte
	}{
		{"Hello, World", []byte("Hello, World\x00")},
		{"", []byte("\x00")},
		{"\x00", []byte("\x00\x00")},
	}

	for i, c := range cases {
		cstr := util.ToCString(c.Given)
		bytes := unsafe.Slice(cstr, len(c.Given)+1)

		if slices.Compare(bytes, c.Wanted) != 0 {
			t.Fatalf("string #%d %q did not match after conversion %q", i, c.Given, string(c.Wanted))
		}
	}
}

func TestToGoString(t *testing.T) {
	cases := []struct {
		Given  []byte
		Wanted string
	}{
		{[]byte("Hello, World\x00"), "Hello, World"},
		{[]byte("\x00"), ""},
		{[]byte("\x00\x00"), ""},
	}

	for i, c := range cases {
		cstr := util.ToGoString(&c.Given[0])
		if cstr != c.Wanted {
			t.Fatalf("string #%d %q did not match after conversion %q", i, string(c.Given), c.Wanted)
		}
	}
}
