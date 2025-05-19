package rgfw

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

type MouseNotifyCallback uintptr

func SetMouseNotify(callback func(win *Window, point Point, status bool)) {
	cb := purego.NewCallback(func(win *Window, point uintptr, status Bool) uintptr {
		p := *(*Point)(unsafe.Pointer(&point))
		callback(win, p, status.IsTrue())
		return 0
	})

	rgfwSetMouseNotifyCallback(MouseNotifyCallback(cb))
}
