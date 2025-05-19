package embedded

import (
	_ "embed"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	dir := filepath.Join(os.TempDir(), "RGFW")
	file := filepath.Join(dir, name)

	if fi, err := os.Stat(file); err == nil {
		if fi.Size() != int64(len(lib)) {
			if err := os.Remove(file); err != nil {
				panic(err)
			}
			if err := os.WriteFile(file, lib, os.ModePerm); err != nil { //nolint:gosec
				panic(err)
			}
		}
	} else {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil { //nolint:gosec
			panic(err)
		}
		if err := os.WriteFile(file, lib, os.ModePerm); err != nil { //nolint:gosec
			panic(err)
		}
	}

	if runtime.GOOS == "windows" {
		if err := os.Setenv("PATH", dir+";"+os.Getenv("PATH")); err != nil {
			panic(err)
		}
	} else {
		if err := os.Setenv("RGFW_PATH", dir); err != nil {
			panic(err)
		}
	}
}
