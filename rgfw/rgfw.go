//go:build purego
// +build purego

package rgfw

import (
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/millken/go/rgfw/embedded"
)

type Point struct {
	X int32
	Y int32
}

type Rect struct {
	X int32
	Y int32
	W int32
	H int32
}

// 修改cRect结构体确保与C语言定义的RGFW_rect完全匹配
type cRect struct {
	x int32
	y int32
	w int32
	h int32
}

func (r *Rect) toC() *cRect {
	if r == nil {
		return nil
	}
	return &cRect{
		x: r.X,
		y: r.Y,
		w: r.W,
		h: r.H,
	}
}

type Area struct {
	W uint32
	H uint32
}

type Monitor struct {
	Name   string
	Rect   Rect
	ScaleX float32
	ScaleY float32
	PhysW  float32
	PhysH  float32
}

type Bool uint8

func toBool(b bool) Bool {
	if b {
		return 1
	}
	return 0
}

func (b Bool) IsTrue() bool {
	return b == 1
}
func (b Bool) IsFalse() bool {
	return b == 0
}

type cEvent struct {
	Type uint8
	// _           [3]byte // padding
	Point   Point
	Vector  Point
	ScaleX  float32
	ScaleY  float32
	Key     uint8
	KeyChar uint8
	Repeat  uint8
	KeyMod  uint8
	Button  uint8
	// _           [7]byte // padding to align double
	Scroll      float64
	Gamepad     uint16
	AxisesCount uint8
	WhichAxis   uint8
	// _           [4]byte // padding to align next field (Axis) to 8 bytes

	Axis [4]Point
	// droppedFiles 和 _win 需要特殊处理
	DroppedFiles      uintptr // char**
	DroppedFilesCount uintptr // size_t
	Win               uintptr // void*
}

type Event struct {
	Type           EventType // which event has been sent?
	Point          Point     // mouse x, y of event (or drop point)
	Vector         Point     // raw mouse movement
	ScaleX, ScaleY float32   // DPI scaling

	Key     Key   // the physical key of the event
	KeyChar uint8 // mapped key char of the event

	Repeat bool // key press event repeated (the key is being held)
	KeyMod KeyMod

	Button uint8   // which mouse (or gamepad) button was pressed
	Scroll float64 // the raw mouse scroll value

	Gamepad     uint16 // which gamepad this event applies to (if applicable to any)
	AxisesCount uint8  // number of axises

	WhichAxis uint8    // which axis was effected
	Axis      [4]Point // x, y of axises (-100 to 100)

	// drag and drop data
	DroppedFiles      []string // dropped files
	DroppedFilesCount int      // how many files were dropped

	Win interface{} // the window this event applies to (for event queue events)
}

type Mouse struct{}

func (m *Mouse) Free() {
	rgfwFreeMouse(m)
}

type Window struct{}

func (w *Window) MakeCurrent() {
	rgfwWindowMakeCurrent(w)
}

func (w *Window) SetIcon(icon []byte, area Area, channels int32) bool {
	return rgfwWindowSetIcon(w, (*byte)(&icon[0]), &area, channels)
}

func (w *Window) SetMouseStandard(mouse MouseIcon) bool {
	return rgfwWindowSetMouseStandard(w, mouse)
}

func (w *Window) SetShouldClose(shouldClose bool) bool {
	return rgfwWindowSetShouldClose(w, toBool(shouldClose))
}

func (w *Window) ShouldClose() bool {
	return rgfwWindowShouldClose(w)
}

func (w *Window) WriteClipboard(text string) {
	rgfwWriteClipboard(ToBytePtr(text), uint32(len(text)+1))
}

func (w *Window) PollEvent() *Event {
	if w == nil {
		return nil
	}
	ptr := rgfwWindowCheckEvent(w)
	if ptr == nil {
		return nil
	}
	ce := (*cEvent)(unsafe.Pointer(ptr))
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

func (w *Window) SwapBuffers() {
	rgfwWindowSwapBuffers(w)
}

func (w *Window) IsRelease(key Key) bool {
	return rgfwIsReleased(w, key)
}

func (w *Window) Close() {
	rgfwWindowClose(w)
}

func (w *Window) SetMouse(mouse *Mouse) {
	rgfwWindowSetMouse(w, mouse)
}

func (w *Window) SetMouseDefault() bool {
	return rgfwWindowSetMouseDefault(w)
}

func (w *Window) ShowMouse(show bool) {
	rgfwWindowShowMouse(w, toBool(show))
}

// Global once to load native library symbols.
var loadOnce sync.Once
var (
	//RGFWDEF void RGFW_setClassName(const char* name);
	rgfwSetClassName func(*byte)
	// RGFWDEF RGFW_window* RGFW_createWindow(
	// 	const char* name, /* name of the window */
	// 	RGFW_rect rect, /* rect of window */
	// 	RGFW_windowFlags flags /* extra arguments ((u32)0 means no flags used)*/
	// ); /*!< function to create a window and struct */

	rgfwCreateWindow func(uintptr, uintptr, WindowFlags) uintptr
	//RGFWDEF void RGFW_window_makeCurrent(RGFW_window* win);
	rgfwWindowMakeCurrent func(*Window)
	//RGFWDEF RGFW_mouse* RGFW_loadMouse(u8* icon, RGFW_area a, i32 channels);
	rgfwLoadMouse func(*byte, *Area, int32) *Mouse
	// RGFWDEF RGFW_bool RGFW_window_setIcon(RGFW_window* win, /*!< source window */
	// 	u8* icon /*!< icon bitmap */,
	// 	RGFW_area a /*!< width and height of the bitmap */,
	// 	i32 channels /*!< how many channels the bitmap has (rgb : 3, rgba : 4) */
	// ); /*!< image MAY be resized by default, set both the taskbar and window icon */
	rgfwWindowSetIcon func(*Window, *byte, *Area, int32) bool
	//RGFWDEF	RGFW_bool RGFW_window_setMouseStandard(RGFW_window* win, u8 mouse);
	rgfwWindowSetMouseStandard func(*Window, MouseIcon) bool
	//RGFWDEF RGFW_bool RGFW_window_shouldClose(RGFW_window* win);
	rgfwWindowShouldClose func(*Window) bool
	//RGFWDEF RGFW_event* RGFW_window_checkEvent(RGFW_window* win);
	rgfwWindowCheckEvent func(*Window) *cEvent
	//RGFWDEF void RGFW_window_setShouldClose(RGFW_window* win, RGFW_bool shouldClose);
	rgfwWindowSetShouldClose func(*Window, Bool) bool
	//RGFWDEF RGFW_bool RGFW_isReleased(RGFW_window* win, RGFW_key key); /*!< if key is released (key code) */
	rgfwIsReleased func(*Window, Key) bool
	//RGFWDEF RGFW_bool RGFW_window_setMouseDefault(RGFW_window* win);
	rgfwWindowSetMouseDefault func(*Window) bool
	//RGFWDEF void RGFW_window_setMouse(RGFW_window* win, RGFW_mouse* mouse);
	rgfwWindowSetMouse func(*Window, *Mouse)
	//RGFWDEF void RGFW_window_showMouse(RGFW_window* win, RGFW_bool show);
	rgfwWindowShowMouse func(*Window, Bool)
	//RGFWDEF void RGFW_writeClipboard(const char* text, u32 textLen); /*!< write text to the clipboard */
	rgfwWriteClipboard func(*byte, uint32)
	//RGFWDEF const char* RGFW_readClipboard(size_t* size);
	rgfwReadClipboard func(*uintptr) *byte
	//RGFWDEF void RGFW_window_swapBuffers(RGFW_window* win); /*!< swap the rendering buffer */
	rgfwWindowSwapBuffers func(*Window)
	//RGFWDEF void RGFW_freeMouse(RGFW_mouse* mouse);
	rgfwFreeMouse func(*Mouse)
	//RGFWDEF void RGFW_window_close(RGFW_window* win); /*!< close the window and free leftover data */
	rgfwWindowClose func(*Window)

	//RGFWDEF double RGFW_getTime(void);
	rgfwGetTime func() float64
	//RGFWDEF u32 RGFW_checkFPS(double startTime, u32 frameCount, u32 fpsCap);
	rgfwCheckFPS func(startTime float64, frameCount uint32, fpsCap uint32) uint32

	/*
		CallBack
	*/
	//RGFWDEF RGFW_mouseNotifyfunc RGFW_setMouseNotifyCallback(RGFW_mouseNotifyfunc func);
	rgfwSetMouseNotifyCallback func(uintptr) uintptr // 修改为使用 uintptr 而不是结构体

	// 添加 MouseNotifyCallback 类型的定义
)

func SetClassName(name string) {
	if len(name) == 0 {
		return
	}

	cStr, ptr := cString(name)
	purego.SyscallN(pSetClassName, ptr)
	runtime.KeepAlive(cStr)

}

func SetClipboard(text string) {
	if len(text) == 0 {
		return
	}

	rgfwWriteClipboard((*byte)(ToBytePtr(text)), uint32(len(text)+1))
}

func GetClipboard() string {

	ptr := rgfwReadClipboard(nil)
	if ptr == nil {
		return ""
	}
	return ToString(ptr)
}

func LoadMouse(icon []byte, area Area, channels int32) *Mouse {
	return rgfwLoadMouse((*byte)(&icon[0]), &area, channels)
}

func GetTime() float64 {
	return rgfwGetTime()
}
func CheckFPS(startTime float64, frameCount uint32, fpsCap uint32) uint32 {
	return rgfwCheckFPS(startTime, frameCount, fpsCap)
}

func CreateWindow(title string, r Rect, flags WindowFlags) *Window {
	// 确保矩形的宽度和高度不为零
	if r.W <= 0 {
		r.W = 640 // 默认宽度
	}
	if r.H <= 0 {
		r.H = 480 // 默认高度
	}

	// 分配在堆上的内存，不会被移动
	rectCopy := r.toC() // 创建副本，避免修改原始结构

	// 创建一个长生命周期的字符串
	cStr, ptr := cString(title)

	// 创建一个引用数组，防止GC回收
	refs := []interface{}{cStr, rectCopy}
	runtime.KeepAlive(refs)

	r1, _, r3 := purego.SyscallN(pCreate, ptr, uintptr(unsafe.Pointer(rectCopy)), uintptr(flags))
	if r3 == 0 {
		panic("rgfw: failed to create window")
	}
	return (*Window)(unsafe.Pointer(r1))
}

var (
	pCreate, pSetClassName uintptr
)

func init() {
	loadOnce.Do(func() {
		if err := embedded.Init(); err != nil {
			panic("rgfw: failed to initialize embedded resources: " + err.Error())
		}
		libHandle, err := loadLibrary(libraryPath())
		if err != nil {
			panic("rgfw: failed to load native library: " + err.Error())
		}
		if libHandle == 0 {
			panic("rgfw: native library not loaded")
		}
		pCreate = loadSymbol(libHandle, "RGFW_createWindow")
		pSetClassName = loadSymbol(libHandle, "RGFW_setClassName")
		// purego.RegisterLibFunc(&rgfwSetClassName, libHandle, "RGFW_setClassName")
		// purego.RegisterLibFunc(&rgfwLoadMouse, libHandle, "RGFW_loadMouse")
		// purego.RegisterLibFunc(&rgfwCreateWindow, libHandle, "RGFW_createWindow")
		// purego.RegisterLibFunc(&rgfwWindowMakeCurrent, libHandle, "RGFW_window_makeCurrent")
		// purego.RegisterLibFunc(&rgfwWindowSetIcon, libHandle, "RGFW_window_setIcon")
		// purego.RegisterLibFunc(&rgfwWindowSetMouseStandard, libHandle, "RGFW_window_setMouseStandard")
		// purego.RegisterLibFunc(&rgfwWindowShouldClose, libHandle, "RGFW_window_shouldClose")
		// purego.RegisterLibFunc(&rgfwWindowCheckEvent, libHandle, "RGFW_window_checkEvent")
		// purego.RegisterLibFunc(&rgfwWindowSetShouldClose, libHandle, "RGFW_window_setShouldClose")
		// purego.RegisterLibFunc(&rgfwIsReleased, libHandle, "RGFW_isReleased")
		// purego.RegisterLibFunc(&rgfwWindowSetMouseDefault, libHandle, "RGFW_window_setMouseDefault")
		// purego.RegisterLibFunc(&rgfwWindowSetMouse, libHandle, "RGFW_window_setMouse")
		// purego.RegisterLibFunc(&rgfwWindowShowMouse, libHandle, "RGFW_window_showMouse")
		// purego.RegisterLibFunc(&rgfwWriteClipboard, libHandle, "RGFW_writeClipboard")
		// purego.RegisterLibFunc(&rgfwReadClipboard, libHandle, "RGFW_readClipboard")
		// purego.RegisterLibFunc(&rgfwWindowSwapBuffers, libHandle, "RGFW_window_swapBuffers")
		// purego.RegisterLibFunc(&rgfwFreeMouse, libHandle, "RGFW_freeMouse")
		// purego.RegisterLibFunc(&rgfwWindowClose, libHandle, "RGFW_window_close")

		// purego.RegisterLibFunc(&rgfwGetTime, libHandle, "RGFW_getTime")
		// purego.RegisterLibFunc(&rgfwCheckFPS, libHandle, "RGFW_checkFPS")

		// purego.RegisterLibFunc(&rgfwSetMouseNotifyCallback, libHandle, "RGFW_setMouseNotifyCallback")
	})
}

// ToBytePtr converts a Go string to a null-terminated C-style string by just appending a null byte,
// if s doesn't already contain one.
func ToBytePtr(s string) *byte {
	size := len(s) + 1
	if index := strings.IndexByte(s, 0); index != -1 {
		size = index + 1
	}

	result := make([]byte, size)
	copy(result, s)
	return &result[0]
}

func cString(s string) ([]byte, uintptr) {
	if s == "" {
		empty := []byte{0}
		return empty, uintptr(unsafe.Pointer(&empty[0]))
	}

	bytes := make([]byte, len(s)+1)
	copy(bytes, s)
	bytes[len(s)] = 0 // 确保以null结尾

	return bytes, uintptr(unsafe.Pointer(&bytes[0]))
}

// ToString converts a null-terminated C-style string into a Go string.
func ToString(p *byte) string {
	if p == nil {
		return ""
	}
	i := 0
	for ptr := unsafe.Pointer(p); *(*byte)(unsafe.Add(ptr, i)) != 0; i++ {
	}
	return string(unsafe.Slice(p, i))
}

// ToStringSlice converts a null-terminated list of C-style strings to a slice of Go strings.
func ToStringSlice(pointers **byte) []string {
	if pointers == nil {
		return nil
	}

	strings := make([]string, 0)

	for ptr := unsafe.Pointer(pointers); *(**byte)(ptr) != nil; ptr = unsafe.Add(ptr, 8) {
		strings = append(strings, ToString(*(**byte)(ptr)))
	}

	return strings
}
