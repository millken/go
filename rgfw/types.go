package rgfw

import "unsafe"

//go:generate stringer -linecomment -output=stringers.go -type=EventType,Format,Icon,WindowFlags

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
type Event struct {
	Type         EventType // which event has been sent?
	Common       CommonEvent
	MouseButton  MouseButtonEvent
	MouseScroll  MouseScrollEvent
	MousePos     MousePosEvent
	Key          KeyEvent
	DataDrop     DataDropEvent
	DataDrag     DataDragEvent
	ScaleUpdated ScaleUpdatedEvent
	// RGFW_commonEvent common; /*!< common event data (e.g.) type and win */
	// RGFW_mouseButtonEvent button; /*!< data for a button press/release */
	// RGFW_mouseScrollEvent scroll; /*!< data for a mouse scroll */
	// RGFW_mousePosEvent mouse; /*!< data for mouse motion events */
	// RGFW_keyEvent key; /*!< data for key press/release/hold events */
	// RGFW_dataDropEvent drop; /*!< dropping a file events */
	// RGFW_dataDragEvent drag; /* data for dragging a file events */
	// RGFW_scaleUpdatedEvent scale; /* data for monitor scaling events */
}

type WindowFlags uint64

const (
	WindowNoBorder       WindowFlags = 1 << 0 // the window doesn't have a border
	WindowNoResize       WindowFlags = 1 << 1 // the window cannot be resized by the user
	WindowAllowDND       WindowFlags = 1 << 2 // the window supports drag and drop
	WindowHideMouse      WindowFlags = 1 << 3 // the window should hide the mouse
	WindowFullscreen     WindowFlags = 1 << 4 // the window is fullscreen by default
	WindowTransparent    WindowFlags = 1 << 5 // the window is transparent
	WindowCenter         WindowFlags = 1 << 6 // center the window on the screen
	WindowScaleToMonitor WindowFlags = 1 << 8 // scale the window to the screen
	WindowHide           WindowFlags = 1 << 9 // the window is hidden
	WindowMaximize       WindowFlags = 1 << 10
	WindowCenterCursor   WindowFlags = 1 << 11
	WindowFloating       WindowFlags = 1 << 12 // create a floating window
	WindowFocusOnShow    WindowFlags = 1 << 13 // focus the window when it's shown
	WindowMinimize       WindowFlags = 1 << 14 // minimize the window
	WindowFocus          WindowFlags = 1 << 15 // if the window is in focus
	WindowOpenGL         WindowFlags = 1 << 17 // create an OpenGL context
	WindowEGL            WindowFlags = 1 << 18 // create an EGL context

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

const (
	EventNone EventType = iota // no event has been sent

	EventKeyPressed  /* a key has been pressed */
	EventKeyReleased /*!< a key has been released */
	/*! key event note
	the code of the key pressed is stored in
	RGFW_event.key.value
	!!Keycodes defined at the bottom of the RGFW_HEADER part of this file!!

	while a string version is stored in
	RGFW_event.key.valueString

	RGFW_event.key.mod holds the current mod
	this means if CapsLock, NumLock are active or not
	*/
	EventMouseButtonPressed  /*!< a mouse button has been pressed (left,middle,right) */
	EventMouseButtonReleased /*!< a mouse button has been released (left,middle,right) */
	EventMouseScroll         /*!< a mouse scroll event */
	EventMousePosChanged     /*!< the position of the mouse has been changed */
	/*! mouse event note
	the x and y of the mouse can be found in the vector, RGFW_x, y

	RGFW_event.button.value holds which mouse button was pressed
	*/
	EventWindowMoved   /*!< the window was moved (by the user) */
	EventWindowResized /*!< the window was resized (by the user), [on WASM this means the browser was resized] */
	EventFocusIn       /*!< window is in focus now */
	EventFocusOut      /*!< window is out of focus now */
	EventMouseEnter    /* mouse entered the window */
	EventMouseLeave    /* mouse left the window */
	EventWindowRefresh /* The window content needs to be refreshed */

	/* attribs change event note
	The event data is sent straight to the window structure
	with win->x, win->y, win->w and win->h
	*/
	EventQuit     /*!< the user clicked the quit button */
	EventDataDrop /*!< a file has been dropped into the window */
	EventDataDrag /*!< the start of a drag and drop event, when the file is being dragged */
	/* drop data note
	The x and y coords of the drop are stored in the vector RGFW_x, y

	RGFW_event.drop.count holds how many files were dropped

	This is also the size of the array which stores all the dropped file string,
	RGFW_event.drop.files
	*/
	EventWindowMaximized /*!< the window was maximized */
	EventWindowMinimized /*!< the window was minimized */
	EventWindowRestored  /*!< the window was restored */
	EventScaleUpdated    /*!< content scale factor changed */
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

type Icon uint8

const (
	IconTaskbar Icon = 1 << 0
	IconWindow  Icon = 1 << 1
	IconBoth         = IconTaskbar | IconWindow
)

func goString(c uintptr) string {
	// We take the address and then dereference it to trick go vet from creating a possible misuse of unsafe.Pointer
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&c))
	if ptr == nil {
		return ""
	}
	var length int
	for {
		if *(*byte)(unsafe.Add(ptr, uintptr(length))) == '\x00' {
			break
		}
		length++
	}
	return string(unsafe.Slice((*byte)(ptr), length))
}

// Converted from RGFW.h types (RGFW_event and related structs).
// Prefix RGFW_ removed and C types mapped to Go equivalents.

type CommonEvent struct {
	Type EventType
	Win  uintptr // RGFW_window* -> opaque pointer
}

type MouseButtonEvent struct {
	Type  EventType
	Win   uintptr
	Value uint8
}

type MouseScrollEvent struct {
	Type EventType
	Win  uintptr
	X    float32
	Y    float32
}

type MousePosEvent struct {
	Type EventType
	Win  uintptr
	X    int32
	Y    int32
	VecX float32
	VecY float32
}

type KeyEvent struct {
	Type   EventType
	Win    uintptr
	Value  Key
	Sym    uint8
	Repeat uint8 // RGFW_bool (u8)
	Mod    KeyMod
}

type DataDropEvent struct {
	Type  EventType
	Win   uintptr
	Files uintptr // char** (use uintptr to represent C array)
	Count uintptr // size_t -> uintptr
}

type DataDragEvent struct {
	Type EventType
	Win  uintptr
	X    int32
	Y    int32
}

type ScaleUpdatedEvent struct {
	Type EventType
	Win  uintptr
	X    float32
	Y    float32
}
