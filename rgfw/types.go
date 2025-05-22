package rgfw

import "unsafe"

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

type WindowFlags uint64

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
