//go:build !windows

package rgfw

import "github.com/ebitengine/purego"

func Load(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_LAZY)
}

func Get(lib uintptr, name string) uintptr {
	addr, err := purego.Dlsym(lib, name)
	if err != nil {
		panic(err)
	}
	return addr
}
