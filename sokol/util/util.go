package util

import "C"

import (
	"runtime"
	"unsafe"
)

// CString is a pointer to a null-terminated C string.
// @fix(judah): Use cgo.Handle or runtime.Pinner to avoid heap allocating every c-string
type CString = *byte

func ToGoString(str CString) string {
	if str == nil {
		panic("cannot convert nil CString to string")
	}

	return C.GoStringN((*C.char)(unsafe.Pointer(str)), C.int(CStringLen(str)))
}

func GoStringN(str CString, n int) string {
	return string(unsafe.Slice((*byte)(unsafe.Pointer(str)), n))
}
func ToCString(str string) CString {
	return (CString)(unsafe.Pointer(C.CString(str)))
}

// ToCStringPin returns a pointer to a pinned Go string (must be null-terminated!)
// and a pinner, which must be kept alive until C is done using the string.
// The pinner will automatically unpin the string when it goes out of scope.
// This is useful for passing strings to C functions that expect a null-terminated string.
func ToCStringPin(str string) (CString, *runtime.Pinner) {
	// 追加 '\0' 保证 C 端安全
	b := append([]byte(str), 0)
	var pinner runtime.Pinner
	pinner.Pin(&b[0])
	return (*byte)(unsafe.Pointer(&b[0])), &pinner
}

func CStringLen(str CString) int {
	ptr := uintptr(unsafe.Pointer(str))
	for *(*byte)((unsafe.Pointer)(ptr)) != 0 {
		ptr += 1
	}

	return int(ptr - uintptr(unsafe.Pointer(str)))
}
