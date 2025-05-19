package rgfw

import "C"
import (
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
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

type Area struct {
	W uint32
	H uint32
}

type Monitor struct {
	Name   *C.char
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
	return rgfwWindowSetIcon(w, (*byte)(&icon[0]), area, channels)
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

type WindowFlags uint32

const (
	WindowNoInitAPI       WindowFlags = 1 << 0  // do NOT init an API (including the software rendering buffer)
	WindowNoBorder        WindowFlags = 1 << 1  // the window doesn't have a border
	WindowNoResize        WindowFlags = 1 << 2  // the window cannot be resized by the user
	WindowAllowDND        WindowFlags = 1 << 3  // the window supports drag and drop
	WindowHideMouse       WindowFlags = 1 << 4  // the window should hide the mouse
	WindowFullscreen      WindowFlags = 1 << 5  // the window is fullscreen by default
	WindowTransparent     WindowFlags = 1 << 6  // the window is transparent
	WindowCenter          WindowFlags = 1 << 7  // center the window on the screen
	WindowOpenglSoftware  WindowFlags = 1 << 8  // use OpenGL software rendering
	WindowCocoaCHDirToRes WindowFlags = 1 << 9  // (cocoa only), change directory to resource folder
	WindowScaleToMonitor  WindowFlags = 1 << 10 // scale the window to the screen
	WindowHide            WindowFlags = 1 << 11 // the window is hidden
	WindowMaximize        WindowFlags = 1 << 12
	WindowCenterCursor    WindowFlags = 1 << 13
	WindowFloating        WindowFlags = 1 << 14 // create a floating window
	WindowFreeOnClose     WindowFlags = 1 << 15 // free the window struct when the window is closed
	WindowFocusOnShow     WindowFlags = 1 << 16 // focus the window when it's shown
	WindowMinimize        WindowFlags = 1 << 17 // minimize the window
	WindowFocus           WindowFlags = 1 << 18 // if the window is in focus

	WindowedFullscreen WindowFlags = WindowNoBorder | WindowMaximize
)

type MouseIcon uint8

const (
	MouseNormal MouseIcon = iota
	MouseArrow
	MouseIbeam
	MouseCrosshair
	MousePointingHand
	MouseResizeEW
	MouseResizeNS
	MouseResizeNWSE
	MouseResizeNESW
	MouseResizeAll
	MouseNotAllowed
	// 11~15 are unused, for alignment
	MouseIconFinal MouseIcon = 16 // padding for alignment
)

// EventType represents the type of an event.
type EventType uint8

func (e EventType) String() string {
	switch e {
	case EventNone:
		return "EventNone"
	case EventKeyPressed:
		return "EventKeyPressed"
	case EventKeyReleased:
		return "EventKeyReleased"
	case EventMouseButtonPressed:
		return "EventMouseButtonPressed"
	case EventMouseButtonReleased:
		return "EventMouseButtonReleased"
	case EventMousePosChanged:
		return "EventMousePosChanged"
	case EventGamepadConnected:
		return "EventGamepadConnected"
	case EventGamepadDisconnected:
		return "EventGamepadDisconnected"
	case EventGamepadButtonPressed:
		return "EventGamepadButtonPressed"
	case EventGamepadButtonReleased:
		return "EventGamepadButtonReleased"
	case EventGamepadAxisMove:
		return "EventGamepadAxisMove"
	case EventWindowMoved:
		return "EventWindowMoved"
	case EventWindowResized:
		return "EventWindowResized"
	case EventFocusIn:
		return "EventFocusIn"
	case EventFocusOut:
		return "EventFocusOut"
	case EventMouseEnter:
		return "EventMouseEnter"
	case EventMouseLeave:
		return "EventMouseLeave"
	case EventWindowRefresh:
		return "EventWindowRefresh"
	case EventQuit:
		return "EventQuit"
	case EventDND:
		return "EventDND"
	case EventDNDInit:
		return "EventDNDInit"
	case EventWindowMaximized:
		return "EventWindowMaximized"
	case EventWindowMinimized:
		return "EventWindowMinimized"
	case EventWindowRestored:
		return "EventWindowRestored"
	default:
	}
	panic("unknown event type")
}

const (
	EventNone EventType = iota // no event has been sent

	// Key events
	EventKeyPressed  // a key has been pressed
	EventKeyReleased // a key has been released
	/*
	   key event note:
	   - the code of the key pressed is stored in Event.key
	   - string version is stored in Event.KeyString
	   - Event.keyMod holds the current keyMod (CapsLock, NumLock, etc.)
	*/

	// Mouse events
	EventMouseButtonPressed  // a mouse button has been pressed (left, middle, right)
	EventMouseButtonReleased // a mouse button has been released (left, middle, right)
	EventMousePosChanged     // the position of the mouse has been changed
	/*
	   mouse event note:
	   - x and y of the mouse can be found in Event.point
	   - Event.button holds which mouse button was pressed
	*/

	// Gamepad events
	EventGamepadConnected      // a gamepad was connected
	EventGamepadDisconnected   // a gamepad was disconnected
	EventGamepadButtonPressed  // a gamepad button was pressed
	EventGamepadButtonReleased // a gamepad button was released
	EventGamepadAxisMove       // an axis of a gamepad was moved
	/*
	   gamepad event note:
	   - Event.gamepad holds which gamepad was altered, if any
	   - Event.button holds which gamepad button was pressed
	   - Event.axis holds the data of all the axes
	   - Event.axisesCount says how many axes there are
	*/

	// Window events
	EventWindowMoved   // the window was moved (by the user)
	EventWindowResized // the window was resized (by the user), [on WASM this means the browser was resized]
	EventFocusIn       // window is in focus now
	EventFocusOut      // window is out of focus now
	EventMouseEnter    // mouse entered the window
	EventMouseLeave    // mouse left the window
	EventWindowRefresh // The window content needs to be refreshed

	/*
	   attribs change event note:
	   - The event data is sent straight to the window structure
	   - win.r.x, win.r.y, win.r.w and win.r.h
	*/

	EventQuit // the user clicked the quit button

	// Drag and Drop events
	EventDND     // a file has been dropped into the window
	EventDNDInit // the start of a dnd event, when the place where the file drop is known
	/*
	   dnd data note:
	   - x and y coords of the drop are stored in Event.point
	   - Event.droppedFilesCount holds how many files were dropped
	   - This is also the size of the array which stores all the dropped file strings, Event.droppedFiles
	*/

	// Window state events
	EventWindowMaximized // the window was maximized
	EventWindowMinimized // the window was minimized
	EventWindowRestored  // the window was restored
	EventScaleUpdated    // content scale factor changed
)

type KeyMod uint8

const (
	KeyModCapsLock   KeyMod = 1 << 0
	KeyModNumLock    KeyMod = 1 << 1
	KeyModControl    KeyMod = 1 << 2
	KeyModAlt        KeyMod = 1 << 3
	KeyModShift      KeyMod = 1 << 4
	KeyModSuper      KeyMod = 1 << 5
	KeyModScrollLock KeyMod = 1 << 6
)

type Key uint8

const (
	KeyNULL     Key = 0
	KeyEscape   Key = '\033'
	KeyBacktick Key = '`'
	Key0        Key = '0'
	Key1        Key = '1'
	Key2        Key = '2'
	Key3        Key = '3'
	Key4        Key = '4'
	Key5        Key = '5'
	Key6        Key = '6'
	Key7        Key = '7'
	Key8        Key = '8'
	Key9        Key = '9'

	KeyMinus     Key = '-'
	KeyEquals    Key = '='
	KeyBackSpace Key = '\b'
	KeyTab       Key = '\t'
	KeySpace     Key = ' '

	KeyA Key = 'a'
	KeyB Key = 'b'
	KeyC Key = 'c'
	KeyD Key = 'd'
	KeyE Key = 'e'
	KeyF Key = 'f'
	KeyG Key = 'g'
	KeyH Key = 'h'
	KeyI Key = 'i'
	KeyJ Key = 'j'
	KeyK Key = 'k'
	KeyL Key = 'l'
	KeyM Key = 'm'
	KeyN Key = 'n'
	KeyO Key = 'o'
	KeyP Key = 'p'
	KeyQ Key = 'q'
	KeyR Key = 'r'
	KeyS Key = 's'
	KeyT Key = 't'
	KeyU Key = 'u'
	KeyV Key = 'v'
	KeyW Key = 'w'
	KeyX Key = 'x'
	KeyY Key = 'y'
	KeyZ Key = 'z'

	KeyPeriod       Key = '.'
	KeyComma        Key = ','
	KeySlash        Key = '/'
	KeyBracket      Key = '{'
	KeyCloseBracket Key = '}'
	KeySemicolon    Key = ';'
	KeyApostrophe   Key = '\''
	KeyBackSlash    Key = '\\'
	KeyReturn       Key = '\n'

	KeyDelete Key = 127 // '\177'
)
const (
	// 自动递增部分，起始值需大于127
	KeyF1 Key = iota + 128
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12

	KeyCapsLock
	KeyShiftL
	KeyControlL
	KeyAltL
	KeySuperL
	KeyShiftR
	KeyControlR
	KeyAltR
	KeySuperR
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyInsert
	KeyEnd
	KeyHome
	KeyPageUp
	KeyPageDown

	KeyNumLock
	KeyKP_Slash
	KeyMultiply
	KeyKP_Minus
	KeyKP_1
	KeyKP_2
	KeyKP_3
	KeyKP_4
	KeyKP_5
	KeyKP_6
	KeyKP_7
	KeyKP_8
	KeyKP_9
	KeyKP_0
	KeyKP_Period
	KeyKP_Return
	KeyScrollLock

	KeyLast = 256
)

// Global once to load native library symbols.
var loadOnce sync.Once
var (
	//RGFWDEF void RGFW_setClassName(const char* name);
	rgfwSetClassName func(*byte)
	rgfwCreateWindow func(string, Rect, WindowFlags) *Window
	//RGFWDEF void RGFW_window_makeCurrent(RGFW_window* win);
	rgfwWindowMakeCurrent func(*Window)
	//RGFWDEF RGFW_mouse* RGFW_loadMouse(u8* icon, RGFW_area a, i32 channels);
	rgfwLoadMouse func(*byte, Area, int32) *Mouse
	// RGFWDEF RGFW_bool RGFW_window_setIcon(RGFW_window* win, /*!< source window */
	// 	u8* icon /*!< icon bitmap */,
	// 	RGFW_area a /*!< width and height of the bitmap */,
	// 	i32 channels /*!< how many channels the bitmap has (rgb : 3, rgba : 4) */
	// ); /*!< image MAY be resized by default, set both the taskbar and window icon */
	rgfwWindowSetIcon func(*Window, *byte, Area, int32) bool
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
	rgfwSetMouseNotifyCallback func(MouseNotifyCallback) // RGFWDEF void RGFW_setMouseNotifyCallback(RGFW_mouseNotifyCallback callback);
)

func SetClassName(name string) {
	if len(name) == 0 {
		return
	}
	loadOnce.Do(func() { registerLibFunc() })
	rgfwSetClassName((*byte)(ToBytePtr(name)))
}

func SetClipboard(text string) {
	if len(text) == 0 {
		return
	}
	loadOnce.Do(func() { registerLibFunc() })
	rgfwWriteClipboard((*byte)(ToBytePtr(text)), uint32(len(text)+1))
}

func GetClipboard() string {
	loadOnce.Do(func() { registerLibFunc() })
	ptr := rgfwReadClipboard(nil)
	if ptr == nil {
		return ""
	}
	return ToString(ptr)
}

func LoadMouse(icon []byte, area Area, channels int32) *Mouse {
	loadOnce.Do(func() { registerLibFunc() })
	return rgfwLoadMouse((*byte)(&icon[0]), area, channels)
}

func GetTime() float64 {
	loadOnce.Do(func() { registerLibFunc() })
	return rgfwGetTime()
}
func CheckFPS(startTime float64, frameCount uint32, fpsCap uint32) uint32 {
	loadOnce.Do(func() { registerLibFunc() })
	return rgfwCheckFPS(startTime, frameCount, fpsCap)
}

func CreateWindow(title string, r Rect, flags WindowFlags) *Window {
	loadOnce.Do(func() { registerLibFunc() })
	return rgfwCreateWindow(title, r, flags)
}

func registerLibFunc() {
	libHandle, err := loadLibrary(libraryPath())
	if err != nil {
		panic("rgfw: failed to load native library: " + err.Error())
	}
	if libHandle == 0 {
		panic("rgfw: native library not loaded")
	}
	purego.RegisterLibFunc(&rgfwSetClassName, libHandle, "RGFW_setClassName")
	purego.RegisterLibFunc(&rgfwLoadMouse, libHandle, "RGFW_loadMouse")
	purego.RegisterLibFunc(&rgfwCreateWindow, libHandle, "RGFW_createWindow")
	purego.RegisterLibFunc(&rgfwWindowMakeCurrent, libHandle, "RGFW_window_makeCurrent")
	purego.RegisterLibFunc(&rgfwWindowSetIcon, libHandle, "RGFW_window_setIcon")
	purego.RegisterLibFunc(&rgfwWindowSetMouseStandard, libHandle, "RGFW_window_setMouseStandard")
	purego.RegisterLibFunc(&rgfwWindowShouldClose, libHandle, "RGFW_window_shouldClose")
	purego.RegisterLibFunc(&rgfwWindowCheckEvent, libHandle, "RGFW_window_checkEvent")
	purego.RegisterLibFunc(&rgfwWindowSetShouldClose, libHandle, "RGFW_window_setShouldClose")
	purego.RegisterLibFunc(&rgfwIsReleased, libHandle, "RGFW_isReleased")
	purego.RegisterLibFunc(&rgfwWindowSetMouseDefault, libHandle, "RGFW_window_setMouseDefault")
	purego.RegisterLibFunc(&rgfwWindowSetMouse, libHandle, "RGFW_window_setMouse")
	purego.RegisterLibFunc(&rgfwWindowShowMouse, libHandle, "RGFW_window_showMouse")
	purego.RegisterLibFunc(&rgfwWriteClipboard, libHandle, "RGFW_writeClipboard")
	purego.RegisterLibFunc(&rgfwReadClipboard, libHandle, "RGFW_readClipboard")
	purego.RegisterLibFunc(&rgfwWindowSwapBuffers, libHandle, "RGFW_window_swapBuffers")
	purego.RegisterLibFunc(&rgfwFreeMouse, libHandle, "RGFW_freeMouse")
	purego.RegisterLibFunc(&rgfwWindowClose, libHandle, "RGFW_window_close")

	purego.RegisterLibFunc(&rgfwGetTime, libHandle, "RGFW_getTime")
	purego.RegisterLibFunc(&rgfwCheckFPS, libHandle, "RGFW_checkFPS")

	purego.RegisterLibFunc(&rgfwSetMouseNotifyCallback, libHandle, "RGFW_setMouseNotifyCallback")
}

func init() {
	runtime.LockOSThread()

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
