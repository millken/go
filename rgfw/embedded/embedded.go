package embedded

import (
	_ "embed"
	"os"
	"path/filepath"
	"runtime"
)

func Init() error {
	if os.Getenv("RGFW_PATH") != "" {
		return nil
	}

	dir := filepath.Join(os.TempDir(), "RGFW")
	file := filepath.Join(dir, name)

	if fi, err := os.Stat(file); err == nil {
		if fi.Size() != int64(len(lib)) {
			if err := os.Remove(file); err != nil {
				return err
			}
			if err := os.WriteFile(file, lib, os.ModePerm); err != nil { //nolint:gosec
				return err
			}
		}
	} else {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil { //nolint:gosec
			return err
		}
		if err := os.WriteFile(file, lib, os.ModePerm); err != nil { //nolint:gosec
			return err
		}
	}

	if runtime.GOOS == "windows" {
		if err := os.Setenv("PATH", dir+";"+os.Getenv("PATH")); err != nil {
			return err
		}
	} else {
		if err := os.Setenv("RGFW_PATH", dir); err != nil {
			return err
		}
	}
	return nil
}
