//go:build !purego
// +build !purego

package rgfw

/*
 #cgo CFLAGS:  -Wno-unused-function -Wno-implicit-function-declaration -I.
 #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/_libs/ -lRGFW_linux_amd64
 #cgo linux LDFLAGS: -lXrandr -lX11 -ldl -lGL
 #cgo darwin,arm64 LDFLAGS: -mmacosx-version-min=11.0 -framework CoreVideo -framework Cocoa -framework OpenGL -framework IOKit -L${SRCDIR}/_libs/ -lRGFW_darwin_arm64

#define RGFW_OPENGL
 #include <RGFW.h>
// Accessor helpers for union fields - cgo can sometimes not expose union members reliably,
// so provide stable C functions to read the fields we need from Go.
static inline RGFW_window* RGFW_event_common_win(const RGFW_event* e) { return e->common.win; }
static inline RGFW_key RGFW_event_key_value(const RGFW_event* e) { return e->key.value; }
static inline unsigned char RGFW_event_key_sym(const RGFW_event* e) { return e->key.sym; }
static inline RGFW_bool RGFW_event_key_repeat(const RGFW_event* e) { return e->key.repeat; }
static inline RGFW_keymod RGFW_event_key_mod(const RGFW_event* e) { return e->key.mod; }
static inline unsigned char RGFW_event_button_value(const RGFW_event* e) { return e->button.value; }
static inline i32 RGFW_event_mouse_x(const RGFW_event* e) { return e->mouse.x; }
static inline i32 RGFW_event_mouse_y(const RGFW_event* e) { return e->mouse.y; }
static inline float RGFW_event_mouse_vecX(const RGFW_event* e) { return e->mouse.vecX; }
static inline float RGFW_event_mouse_vecY(const RGFW_event* e) { return e->mouse.vecY; }
static inline float RGFW_event_scroll_x(const RGFW_event* e) { return e->scroll.x; }
static inline float RGFW_event_scroll_y(const RGFW_event* e) { return e->scroll.y; }
static inline size_t RGFW_event_drop_count(const RGFW_event* e) { return e->drop.count; }
static inline char** RGFW_event_drop_files(const RGFW_event* e) { return e->drop.files; }
static inline float RGFW_event_scale_x(const RGFW_event* e) { return e->scale.x; }
static inline float RGFW_event_scale_y(const RGFW_event* e) { return e->scale.y; }
extern void goMouseCallback(RGFW_window* win, i32 x, i32 y, RGFW_bool status);
*/
import "C"
import (
	"unsafe"
)

type Format uint8

const (
	FormatRGB8 Format = iota
	FormatBGR8
	FormatRGBA8
	FormatARGB8
	FormatBGRA8
	FormatABGR8
	FormatCount
)

// RGFWDEF void RGFW_setClassName(const char* name);
func SetClassName(name string) {
	cStr := C.CString(name)
	defer C.free(unsafe.Pointer(cStr))
	C.RGFW_setClassName(cStr)
}

// RGFWDEF RGFW_mouse* RGFW_loadMouse(u8* data, i32 w, i32 h, RGFW_format format);
func LoadMouse(icon []uint8, w, h int32, format Format) *Mouse {
	cIcon := C.CBytes(icon)
	defer C.free(cIcon)
	ptr := C.RGFW_loadMouse((*C.u8)(cIcon), C.i32(w), C.i32(h), C.RGFW_format(format))
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

// RGFWDEF RGFW_bool RGFW_window_setIconEx(RGFW_window* win, u8* data, i32 w, i32 h, RGFW_format format, RGFW_icon type);
func (w *Window) SetIconEx(data []byte, width, height int32, format Format, icon Icon) bool {
	if len(data) == 0 {
		return false
	}
	cIcon := C.CBytes(data)
	defer C.free(cIcon)
	return C.RGFW_window_setIconEx(w.ptr, (*C.u8)(cIcon), C.i32(width), C.i32(height), C.RGFW_format(format), C.RGFW_icon(icon)) != 0
}

// RGFWDEF void RGFW_pollEvents(void);
func PollEvents() {
	C.RGFW_pollEvents()
}

// 声明一个全局变量来存储Go回调函数
var mouseNotifyCallback func(win *Window, point Point, status bool)

//export goMouseCallback
func goMouseCallback(win *C.RGFW_window, x, y C.i32, status C.RGFW_bool) {
	if mouseNotifyCallback != nil {
		mouseNotifyCallback(&Window{ptr: win}, Point{X: int32(x), Y: int32(y)}, status != 0)
	}
}

// RGFWDEF RGFW_mouseNotifyfunc RGFW_setMouseNotifyCallback(RGFW_mouseNotifyfunc func);
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

// RGFWDEF RGFW_window* RGFW_createWindow(const char* name, i32 x, i32 y, i32 w, i32 h,  RGFW_windowFlags flags);
func CreateWindow(name string, rect Rect, flags WindowFlags) *Window {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cFlags := C.uint64_t(flags)
	ptr := C.RGFW_createWindow(cName, C.i32(rect.X), C.i32(rect.Y), C.i32(rect.W), C.i32(rect.H), C.RGFW_windowFlags(cFlags))
	if ptr == nil {
		return nil
	}
	return &Window{ptr}
}

func (w *Window) MakeCurrentOpenGL() {
	C.RGFW_window_makeCurrentWindow_OpenGL(w.ptr)
}

func (w *Window) ShouldClose() bool {
	return C.RGFW_window_shouldClose(w.ptr) != 0
}

func (w *Window) SetExitKey(key Key) {
	C.RGFW_window_setExitKey(w.ptr, C.RGFW_key(key))
}

func (w *Window) SetShouldClose(shouldClose bool) {
	if shouldClose {
		C.RGFW_window_setShouldClose(w.ptr, 1)
	} else {
		C.RGFW_window_setShouldClose(w.ptr, 0)
	}
}

// RGFWDEF RGFW_bool RGFW_window_setIcon(RGFW_window* win, u8* data, i32 w, i32 h, RGFW_format format);
func (w *Window) SetIcon(icon []uint8, width, height int32, format Format) bool {
	if len(icon) == 0 {
		return false
	}
	cIcon := C.CBytes(icon)
	defer C.free(cIcon)
	return C.RGFW_window_setIcon(w.ptr, (*C.u8)(cIcon), C.i32(width), C.i32(height), C.RGFW_format(format)) != 0
}

// RGFWDEF RGFW_bool RGFW_isReleased(RGFW_window* win, RGFW_key key);
func (w *Window) IsRelease(key Key) bool {
	return C.RGFW_window_isKeyReleased(w.ptr, C.RGFW_key(key)) != 0
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
func (w *Window) SwapBuffersOpenGL() {
	C.RGFW_window_swapBuffers_OpenGL(w.ptr)
}

// RGFWDEF void RGFW_window_swapInterval_OpenGL(RGFW_window* win, i32 swapInterval);
func (w *Window) SwapIntervalOpenGL(interval int32) {
	C.RGFW_window_swapInterval_OpenGL(w.ptr, C.i32(interval))
}

func (w *Window) CheckEvent(e *Event) bool {
	if e == nil {
		return false
	}

	var cev C.RGFW_event
	ok := C.RGFW_window_checkEvent(w.ptr, &cev)

	// common: set event type (cev.type uses the C field name `type` which is a Go keyword,
	// so read the first byte of the C union directly)
	e.Type = EventType(*(*uint8)(unsafe.Pointer(&cev)))

	// map common win pointer using C accessor
	cw := C.RGFW_event_common_win(&cev)
	if cw != nil {
		e.Common.Win = uintptr(unsafe.Pointer(cw))
	} else {
		e.Common.Win = 0
	}

	switch e.Type {
	case EventKeyPressed, EventKeyReleased:
		e.Key.Value = Key(C.RGFW_event_key_value(&cev))
		e.Key.Sym = uint8(C.RGFW_event_key_sym(&cev))
		e.Key.Repeat = uint8(C.RGFW_event_key_repeat(&cev))
		e.Key.Mod = KeyMod(uint8(C.RGFW_event_key_mod(&cev)))
		e.Key.Win = uintptr(unsafe.Pointer(cw))
	case EventMouseButtonPressed, EventMouseButtonReleased:
		e.MouseButton.Value = uint8(C.RGFW_event_button_value(&cev))
		e.MouseButton.Win = uintptr(unsafe.Pointer(cw))
	case EventMousePosChanged:
		e.MousePos.X = int32(C.RGFW_event_mouse_x(&cev))
		e.MousePos.Y = int32(C.RGFW_event_mouse_y(&cev))
		e.MousePos.VecX = float32(C.RGFW_event_mouse_vecX(&cev))
		e.MousePos.VecY = float32(C.RGFW_event_mouse_vecY(&cev))
		e.MousePos.Win = uintptr(unsafe.Pointer(cw))
	case EventMouseScroll:
		e.MouseScroll.X = float32(C.RGFW_event_scroll_x(&cev))
		e.MouseScroll.Y = float32(C.RGFW_event_scroll_y(&cev))
		e.MouseScroll.Win = uintptr(unsafe.Pointer(cw))
	case EventDataDrop, EventDataDrag:
		// convert char** files to []string
		e.DataDrop.Count = uintptr(C.RGFW_event_drop_count(&cev))
		filesPtr := C.RGFW_event_drop_files(&cev)
		if filesPtr != nil && C.RGFW_event_drop_count(&cev) > 0 {
			e.DataDrop.Files = uintptr(unsafe.Pointer(filesPtr))
		} else {
			e.DataDrop.Files = 0
		}
		e.DataDrop.Win = uintptr(unsafe.Pointer(cw))
	case EventScaleUpdated:
		e.ScaleUpdated.X = float32(C.RGFW_event_scale_x(&cev))
		e.ScaleUpdated.Y = float32(C.RGFW_event_scale_y(&cev))
		e.ScaleUpdated.Win = uintptr(unsafe.Pointer(cw))
	default:
		// for other events, nothing else to map for now
	}
	return ok != 0
}

// RGFWDEF void RGFW_window_close(RGFW_window* win);
func (w *Window) Close() {
	if w == nil {
		return
	}

	C.RGFW_window_close(w.ptr)
}
