package picasso

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo darwin CXXFLAGS:  -std=c++11 -I${SRCDIR} -mmacosx-version-min=15.4
#cgo LDFLAGS:  -lpicasso2_sw -lstdc++
#cgo darwin LDFLAGS: -mmacosx-version-min=15.4 -L${SRCDIR}/_libs/darwin/arm64  -framework CoreGraphics -framework CoreText
#include "include/picasso.h"
#include <stdlib.h>
#include <stddef.h>
#include <stdint.h>
#include "cgo_helpers.h"

*/
import "C"
import (
	"unsafe"
)

const (
	True  Bool = 1
	False Bool = 0
)

type WrapperByte struct {
	Size int
	Ptr  *Byte
}

func MakeByte(size int) *WrapperByte {
	if size <= 0 {
		return nil
	}
	b := C.malloc(C.size_t(size))
	if b == nil {
		return nil
	}
	return &WrapperByte{
		Size: size,
		Ptr:  (*Byte)(b),
	}
}

func (w *WrapperByte) Byte() *Byte {
	if w == nil || w.Ptr == nil {
		return nil
	}
	return w.Ptr
}
func (w *WrapperByte) Free() {
	if w == nil || w.Ptr == nil {
		return
	}
	C.free(unsafe.Pointer(w.Ptr))
	w.Ptr = nil
	w.Size = 0
}

func (w *WrapperByte) Get() []byte {
	if w == nil || w.Ptr == nil {
		return nil
	}
	goBuf := unsafe.Slice((*byte)(w.Ptr), w.Size)
	return goBuf
}

func (w *WrapperByte) Set(data []byte) {
	if w == nil || w.Ptr == nil {
		return
	}
	if len(data) > w.Size {
		return
	}
	goBuf := unsafe.Slice((*byte)(w.Ptr), w.Size)
	copy(goBuf, data)
}

func (w *WrapperByte) SetByte(data []byte) {
	if w == nil || w.Ptr == nil {
		return
	}
	if len(data) > w.Size {
		return
	}
	goBuf := unsafe.Slice((*byte)(w.Ptr), w.Size)
	for i := 0; i < len(data); i++ {
		goBuf[i] = data[i]
	}
}

func FromByte(data []byte) *WrapperByte {
	if len(data) <= 0 {
		return nil
	}
	w := MakeByte(len(data))
	if w == nil {
		return nil
	}
	w.Set(data)
	return w
}
