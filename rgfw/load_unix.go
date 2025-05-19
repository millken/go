//go:build darwin || linux
// +build darwin linux

package rgfw

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/ebitengine/purego"
)

func libraryPath() string {
	var name string
	var paths []string

	rgfwPath := os.Getenv("RGFW_PATH")
	execPath, _ := os.Executable()
	dir := filepath.Dir(execPath)

	switch runtime.GOOS {
	case "linux":
		name = "libRGFW.so"
		paths = []string{rgfwPath, dir}
	case "darwin":
		name = "libRGFW.dylib"
		paths = []string{rgfwPath, dir, filepath.Join(dir, "..", "Frameworks")}
	}

	for _, v := range paths {
		n := filepath.Join(v, name)
		if _, err := os.Stat(n); err == nil {
			name = n
			break
		}
	}

	return name
}

func loadLibrary(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
}

func loadSymbol(lib uintptr, name string) uintptr {
	ptr, err := purego.Dlsym(lib, name)
	if err != nil {
		panic("rgfw: failed to load symbol " + name + ": " + err.Error())
	}
	return ptr
}
