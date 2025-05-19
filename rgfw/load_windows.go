package rgfw

import (
	"syscall"
)

func libraryPath() string {
	return "NGFW.dll"
}

func loadLibrary(name string) (uintptr, error) {
	handle, err := syscall.LoadLibrary(name)
	return uintptr(handle), err
}

func loadSymbol(lib uintptr, name string) uintptr {
	ptr, err := syscall.GetProcAddress(syscall.Handle(lib), name)
	if err != nil {
		panic("rgfw: failed to load symbol " + name + ": " + err.Error())
	}
	return ptr
}
