package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"os"
	"runtime"

	"github.com/millken/go/rgfw"
)

var baseIcon = [4 * 3 * 3]uint8{
	0xFF, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0xFF,
}

func getPNGSize(reader io.Reader) (width, height int, err error) {
	config, err := png.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err
	}
	return config.Width, config.Height, nil
}
func main() {
	runtime.LockOSThread()
	rgfw.SetClassName("RGFW Example")
	win := rgfw.CreateWindow("RGFW Example Window", rgfw.Rect{X: 500, Y: 500, W: 500, H: 500}, rgfw.WindowCenter|rgfw.WindowNoResize)
	rgfw.SetMouseNotify(func(win *rgfw.Window, point rgfw.Point, status bool) {
		fmt.Printf("Mouse moved to %d x %d with status %v\n", point.X, point.Y, status)
	})
	icon, err := os.ReadFile("./logo.png")
	if err != nil {
		fmt.Println("Failed to read icon file:", err)
		return
	}
	buffer := bytes.NewReader(icon)
	w, h, _ := getPNGSize(buffer)
	fmt.Printf("内存尺寸: %d x %d\n", w, h)

	win.SetExitKey(rgfw.KeyEscape)
	win.SetIconEx(baseIcon[:], 3, 3, rgfw.FormatRGBA8, rgfw.IconWindow)
	win.SetIconEx(icon[:], int32(w), int32(h), rgfw.FormatRGB8, rgfw.IconTaskbar)
	for !win.ShouldClose() {
		rgfw.PollEvents()
	}
	win.Close()
}
