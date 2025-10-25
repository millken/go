package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/millken/go/rgfw"
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
	runtime.LockOSThread()
	rgfw.SetClassName("RGFW Example")
	win := rgfw.CreateWindow("RGFW Example Window", rgfw.Rect{X: 500, Y: 500, W: 500, H: 500}, rgfw.WindowCenter|rgfw.WindowAllowDND|rgfw.WindowOpenGL)
	rgfw.SetMouseNotify(func(win *rgfw.Window, point rgfw.Point, status bool) {
		fmt.Printf("Mouse moved to %d x %d with status %v\n", point.X, point.Y, status)
	})
	fmt.Println("窗口图标和鼠标指针示例")
	win.MakeCurrentOpenGL()
	fmt.Println("设置窗口图标和鼠标指针")
	win.SetIconEx(icon[:], 3, 3, rgfw.FormatBGRA8, rgfw.IconBoth)
	win.SetMouseStandard(rgfw.MouseResizeNESW)
	mouse := rgfw.LoadMouse(icon[:], 3, 3, rgfw.FormatBGRA8)
	// 初始化 OpenGL 绑定
	if err := gl.Init(); err != nil {
		panic(err)
	}
	win.SetExitKey(rgfw.KeyEscape)
	win.SwapIntervalOpenGL(1)
	fmt.Printf("OpenGL Version: %s\n", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Printf("OpenGL Vendor: %s\n", gl.GoStr(gl.GetString(gl.VENDOR)))
	frames := uint32(0)
	// 按需渲染：交互活跃时 60 FPS，空闲时 10 FPS，活跃窗口 2 秒
	adaptive := NewAdaptiveLimiter(60, 10, 2*time.Second)
	for !win.ShouldClose() {
		var event = rgfw.GetEvent()
		for win.CheckEvent(event) {
			adaptive.MarkActive()
			switch event.Type {
			case rgfw.EventQuit:
				break
			}
			rgfw.ReleaseEvent(event)
		}
		// rgfw.PollEvents()
		if win.IsKeyRelease(rgfw.KeyW) {
			win.SetMouseDefault()
		} else if win.IsKeyRelease(rgfw.KeyE) {
			win.SetMouse(mouse)
		} else if win.IsKeyRelease(rgfw.KeyQ) {
			win.ShowMouse(false)
		} else if win.IsKeyRelease(rgfw.KeyT) {
			win.ShowMouse(true)
		} else if win.IsKeyRelease(rgfw.KeyDown) {
			rgfw.SetClipboard("DOWN 刺猬")
		} else if win.IsKeyRelease(rgfw.KeyUp) {
			text := rgfw.GetClipboard()
			if text != "" {
				fmt.Printf("pasted '%s'\n", text)
			}
		}

		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.Begin(gl.TRIANGLES)
		gl.Color3f(1.0, 0.0, 0.0)
		gl.Vertex2f(-0.6, -0.75)
		gl.Color3f(0.0, 1.0, 0.0)
		gl.Vertex2f(0.6, -0.75)
		gl.Color3f(0.0, 0.0, 1.0)
		gl.Vertex2f(0.0, 0.75)
		gl.End()
		win.SwapBuffersOpenGL()
		frames++
		adaptive.Sleep()
	}
	fmt.Printf("总帧数: %d\n", frames)
	mouse.Destroy()
	win.Close()
}

// NewFPSLimiter 创建一个简单的帧率限制器
// targetFPS <= 0 时退化为不限制
func NewFPSLimiter(targetFPS int) *fpsLimiter {
	if targetFPS <= 0 {
		return &fpsLimiter{frameDur: 0}
	}
	return &fpsLimiter{frameDur: time.Second / time.Duration(targetFPS)}
}

type fpsLimiter struct {
	frameDur time.Duration
	last     time.Time
}

// SetTargetFPS 动态设置目标 FPS
func (f *fpsLimiter) SetTargetFPS(fps int) {
	if fps <= 0 {
		f.frameDur = 0
		f.last = time.Time{}
		return
	}
	dur := time.Second / time.Duration(fps)
	if dur != f.frameDur {
		f.frameDur = dur
		// 重新计时，避免切换 FPS 时产生长短不一的第一帧
		f.last = time.Now()
	}
}

func (f *fpsLimiter) Sleep() {
	if f.frameDur <= 0 {
		return
	}
	now := time.Now()
	if f.last.IsZero() {
		f.last = now
		return
	}
	elapsed := now.Sub(f.last)
	if remaining := f.frameDur - elapsed; remaining > 0 {
		time.Sleep(remaining)
		f.last = f.last.Add(f.frameDur)
	} else {
		// 超时则从现在重新计时，避免长时间落后
		f.last = now
	}
}

// 自适应 FPS 限制器：活跃期 highFPS，空闲期 lowFPS
/*
用户输入、窗口移动/大小变化等事件发生后，接下来的 2 秒内以 60 FPS 渲染，保证交互流畅。
超过 2 秒没有任何事件时自动降到 10 FPS，减少 CPU 占用与功耗。
你可以调整参数满足不同需求，例如调成 NewAdaptiveLimiter(120, 5, 1*time.Second)。
*/
type adaptiveLimiter struct {
	highFPS, lowFPS int
	activeDur       time.Duration
	lastActive      time.Time
	base            *fpsLimiter
}

func NewAdaptiveLimiter(highFPS, lowFPS int, activeDur time.Duration) *adaptiveLimiter {
	al := &adaptiveLimiter{
		highFPS:   highFPS,
		lowFPS:    lowFPS,
		activeDur: activeDur,
		base:      NewFPSLimiter(lowFPS),
	}
	// 启动时视为活跃一小段时间，避免一进程序就很低 FPS
	al.lastActive = time.Now()
	al.base.SetTargetFPS(highFPS)
	return al
}

// MarkActive 标记活跃（例如有输入事件/窗口变化）
func (a *adaptiveLimiter) MarkActive() { a.lastActive = time.Now() }

func (a *adaptiveLimiter) Sleep() {
	// 根据最近活跃时间切换目标 FPS
	if time.Since(a.lastActive) < a.activeDur {
		a.base.SetTargetFPS(a.highFPS)
	} else {
		a.base.SetTargetFPS(a.lowFPS)
	}
	a.base.Sleep()
}
