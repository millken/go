package main

import (
	"fmt"

	"github.com/millken/go/rgfw"
	_ "github.com/millken/go/rgfw/embedded"
)

var icon = [4 * 3 * 3]uint8{
	0xFF, 0x00, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0x00,
	0xFF, 0xFF, 0x00, 0xFF,
	0xFF, 0xFF, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0xFF,
	0xFF, 0x00, 0x00, 0xFF,
}

func main() {
	// os.Setenv("RGFW_PATH", "/Users/millken/github.com/millken/go/rgfw/RGFW/")
	rgfw.SetClassName("RGFW Example")
	win := rgfw.CreateWindow("RGFW Example Window", rgfw.Rect{X: 500, Y: 500, W: 500, H: 500}, rgfw.WindowCenter|rgfw.WindowAllowDND)
	rgfw.SetMouseNotify(func(win *rgfw.Window, point rgfw.Point, status bool) {
		fmt.Printf("Mouse moved to %d x %d with status %v\n", point.X, point.Y, status)
	})
	win.MakeCurrent()
	win.SetIcon(icon[:], rgfw.Area{3, 3}, 4)
	win.SetMouseStandard(rgfw.MouseResizeNESW)
	mouse := rgfw.LoadMouse(icon[:], rgfw.Area{3, 3}, 4)
	frames := uint32(0)
	fps := uint32(0)
	startTime := rgfw.GetTime()
	for !win.ShouldClose() {
		var event *rgfw.Event
		for event = win.PollEvent(); event != nil; event = win.PollEvent() {
			switch event.Type {
			case rgfw.EventQuit:
				win.SetShouldClose(true)
				break
			case rgfw.EventWindowResized:
				if event.Point.X > 0 && event.Point.Y > 0 {
					fmt.Printf("Window resized to %d x %d\n", event.Point.X, event.Point.Y)
				}
			case rgfw.EventMouseButtonPressed:
				fmt.Printf("Mouse button pressed at %d x %d\n", event.Point.X, event.Point.Y)
			case rgfw.EventMouseButtonReleased:
				fmt.Printf("Mouse button released at %d x %d\n", event.Point.X, event.Point.Y)
			case rgfw.EventGamepadButtonPressed:
				fmt.Printf("Gamepad button pressed at %d x %d\n", event.Point.X, event.Point.Y)
			case rgfw.EventGamepadButtonReleased:
				fmt.Printf("Gamepad button released at %d x %d\n", event.Point.X, event.Point.Y)
			}
		}
		if win.IsRelease(rgfw.KeySpace) {
			fmt.Println("fps:", fps)
		} else if win.IsRelease(rgfw.KeyW) {
			win.SetMouseDefault()
		} else if win.IsRelease(rgfw.KeyE) {
			win.SetMouse(mouse)
		} else if win.IsRelease(rgfw.KeyQ) {
			win.ShowMouse(false)
		} else if win.IsRelease(rgfw.KeyT) {
			win.ShowMouse(true)
		} else if win.IsRelease(rgfw.KeyDown) {
			rgfw.SetClipboard("DOWN 刺猬")
		} else if win.IsRelease(rgfw.KeyUp) {
			text := rgfw.GetClipboard()
			if text != "" {
				fmt.Printf("pasted '%s'\n", text)
			}
		}
		win.SwapBuffers()
		fps = rgfw.CheckFPS(startTime, frames, 60)
		frames++

	}

	mouse.Free()
	win.Close()
}
