//go:build !purego
// +build !purego

package rgfw

/*
 #cgo CFLAGS:  -Wno-unused-function -Wno-implicit-function-declaration -I.
 #cgo linux LDFLAGS: -ldl -lpthread -lX11 -lXrandr -lGL
 #cgo darwin,arm64 LDFLAGS: -mmacosx-version-min=11.0 -framework CoreVideo -framework Cocoa -framework OpenGL -framework IOKit -L${SRCDIR}/_libs/ -lRGFW_darwin_arm64

 #include <RGFW.h>

extern void goMouseCallback(RGFW_window* win, RGFW_point point, RGFW_bool status);
*/
import "C"
import (
	"unsafe"
)

// RGFWDEF void RGFW_setClassName(const char* name);
func SetClassName(name string) {
	cStr := C.CString(name)
	defer C.free(unsafe.Pointer(cStr))
	C.RGFW_setClassName(cStr)
}

func GetTime() float64 {
	return float64(C.RGFW_getTime())
}

// RGFWDEF u32 RGFW_checkFPS(double startTime, u32 frameCount, u32 fpsCap);
func CheckFPS(startTime float64, frameCount uint32, fpsCap uint32) uint32 {
	return uint32(C.RGFW_checkFPS(C.double(startTime), C.uint32_t(frameCount), C.uint32_t(fpsCap)))
}

// RGFWDEF RGFW_mouse* RGFW_loadMouse(u8* icon, RGFW_area a, i32 channels);
func LoadMouse(icon []uint8, area Area, channels int32) *Mouse {
	cIcon := C.CBytes(icon)
	defer C.free(cIcon)
	cArea := areaToC(area)
	ptr := C.RGFW_loadMouse((*C.u8)(cIcon), cArea, C.i32(channels))
	if ptr == nil {
		return nil
	}
	return &Mouse{(*C.RGFW_mouse)(ptr)}
}

// RGFWDEF void RGFW_writeClipboard(const char* text, u32 textLen);
func SetClipboard(text string) {
	if len(text) == 0 {
		return
	}

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.RGFW_writeClipboard(cText, C.u32(len(text)))
}

// RGFWDEF const char* RGFW_readClipboard(size_t* size);
func GetClipboard() string {

	ptr := C.RGFW_readClipboard(nil)
	if ptr == nil {
		return ""
	}
	return goString(uintptr(unsafe.Pointer(ptr)))
}

// 声明一个全局变量来存储Go回调函数
var mouseNotifyCallback func(win *Window, point Point, status bool)

//export goMouseCallback
func goMouseCallback(win *C.RGFW_window, point C.RGFW_point, status C.RGFW_bool) {
	if mouseNotifyCallback != nil {
		mouseNotifyCallback(&Window{ptr: win}, Point{X: int32(point.x), Y: int32(point.y)}, status != 0)
	}
}

func SetMouseNotify(callback func(*Window, Point, bool)) {
	mouseNotifyCallback = callback
	C.RGFW_setMouseNotifyCallback(C.RGFW_mouseNotifyfunc(C.goMouseCallback))
}

type Window struct {
	ptr *C.RGFW_window
}

type Mouse struct {
	ptr *C.RGFW_mouse
}

func (m *Mouse) Destroy() {
	if m.ptr != nil {
		C.RGFW_freeMouse(unsafe.Pointer(m.ptr))
		m.ptr = nil
	}
}

func rectToC(rect Rect) C.RGFW_rect {
	return C.RGFW_rect{
		x: C.int(rect.X),
		y: C.int(rect.Y),
		w: C.int(rect.W),
		h: C.int(rect.H),
	}
}

func areaToC(area Area) C.RGFW_area {
	return C.RGFW_area{
		w: C.u32(area.W),
		h: C.u32(area.H),
	}
}

// RGFW_window* RGFW_createWindow(const char* name, RGFW_rect rect, RGFW_windowFlags flags);
func CreateWindow(name string, rect Rect, flags WindowFlags) *Window {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cFlags := C.uint64_t(flags)
	ptr := C.RGFW_createWindow(cName, rectToC(rect), C.RGFW_windowFlags(cFlags))
	if ptr == nil {
		return nil
	}
	return &Window{ptr}
}

func (w *Window) MakeCurrent() {
	C.RGFW_window_makeCurrent(w.ptr)
}

func (w *Window) ShouldClose() bool {
	return C.RGFW_window_shouldClose(w.ptr) != 0
}

func (w *Window) SetShouldClose(shouldClose bool) {
	if shouldClose {
		C.RGFW_window_setShouldClose(w.ptr, 1)
	} else {
		C.RGFW_window_setShouldClose(w.ptr, 0)
	}
}

// RGFW_bool RGFW_window_setIcon(RGFW_window* win, u8* icon, RGFW_area a, i32 channels)
func (w *Window) SetIcon(icon []uint8, area Area, channels int32) bool {
	if len(icon) == 0 {
		return false
	}
	cIcon := C.CBytes(icon)
	defer C.free(cIcon)
	cArea := areaToC(area)
	return C.RGFW_window_setIcon(w.ptr, (*C.u8)(cIcon), cArea, C.i32(channels)) != 0
}

// RGFWDEF RGFW_bool RGFW_isReleased(RGFW_window* win, RGFW_key key);
func (w *Window) IsRelease(key Key) bool {
	return C.RGFW_isReleased(w.ptr, C.RGFW_key(key)) != 0
}

// RGFWDEF RGFW_bool RGFW_window_setMouseDefault(RGFW_window* win);
func (w *Window) SetMouseDefault() bool {
	return C.RGFW_window_setMouseDefault(w.ptr) != 0
}

// RGFWDEF	RGFW_bool RGFW_window_setMouseStandard(RGFW_window* win, u8 mouse);
func (w *Window) SetMouseStandard(mouse MouseIcon) bool {
	return C.RGFW_window_setMouseStandard(w.ptr, C.u8(mouse)) != 0
}

// RGFWDEF void RGFW_window_setMouse(RGFW_window* win, RGFW_mouse* mouse);
func (w *Window) SetMouse(mouse *Mouse) {
	if mouse == nil {
		return
	}
	C.RGFW_window_setMouse(w.ptr, unsafe.Pointer(mouse.ptr))
}

// RGFWDEF void RGFW_window_showMouse(RGFW_window* win, RGFW_bool show);
func (w *Window) ShowMouse(show bool) {
	if show {
		C.RGFW_window_showMouse(w.ptr, 1)
	} else {
		C.RGFW_window_showMouse(w.ptr, 0)
	}
}
func (w *Window) SwapBuffers() {
	C.RGFW_window_swapBuffers(w.ptr)
}

// RGFWDEF void RGFW_window_swapInterval(RGFW_window* win, i32 interval);
func (w *Window) SwapInterval(interval int32) {
	C.RGFW_window_swapInterval(w.ptr, C.i32(interval))
}

func (w *Window) PollEvent() *Event {
	// RGFW_event* RGFW_window_pollEvent(RGFW_window* win);
	eventPtr := C.RGFW_window_checkEvent(w.ptr)
	if eventPtr == nil {
		return nil
	}
	ce := (*cEvent)(unsafe.Pointer(eventPtr))
	// 手动转换为 Go 的 Event
	ev := &Event{
		Type:              EventType(ce.Type),
		Point:             ce.Point,
		Vector:            ce.Vector,
		ScaleX:            ce.ScaleX,
		ScaleY:            ce.ScaleY,
		KeyChar:           ce.KeyChar,
		Repeat:            ce.Repeat == 1,
		KeyMod:            KeyMod(ce.KeyMod),
		Button:            ce.Button,
		Scroll:            ce.Scroll,
		Gamepad:           ce.Gamepad,
		AxisesCount:       ce.AxisesCount,
		WhichAxis:         ce.WhichAxis,
		Axis:              ce.Axis,
		DroppedFilesCount: int(ce.DroppedFilesCount),
		// DroppedFiles 需要额外处理
		// Win 需要额外处理
	}
	return ev
}

// RGFWDEF void RGFW_window_close(RGFW_window* win);
func (w *Window) Close() {
	if w == nil {
		return
	}

	C.RGFW_window_close(w.ptr)
}
