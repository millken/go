package c2go

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// AutoFree 自动释放C分配的内存
type AutoFree struct {
	ptr unsafe.Pointer
}

func NewAutoFree(size C.size_t) *AutoFree {
	return &AutoFree{
		ptr: C.malloc(size),
	}
}

func (a *AutoFree) Ptr() unsafe.Pointer {
	return a.ptr
}

func (a *AutoFree) Free() {
	if a.ptr != nil {
		C.free(a.ptr)
		a.ptr = nil
	}
}

// StringToC 转换Go字符串到C字符串，自动管理内存
func StringToC(s string) *AutoFree {
	cstr := C.CString(s)
	return &AutoFree{ptr: unsafe.Pointer(cstr)}
}

// Handle 资源句柄
type Handle struct {
	ptr      unsafe.Pointer
	freeFunc func(unsafe.Pointer)
}

func NewHandle(ptr unsafe.Pointer, freeFunc func(unsafe.Pointer)) *Handle {
	return &Handle{
		ptr:      ptr,
		freeFunc: freeFunc,
	}
}

func (h *Handle) Ptr() unsafe.Pointer {
	return h.ptr
}

func (h *Handle) Close() {
	if h.ptr != nil && h.freeFunc != nil {
		h.freeFunc(h.ptr)
		h.ptr = nil
	}
}
