package common

import (
	"fmt"
	"syscall"
	"unsafe"
)

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getCursorPos     = user32.NewProc("GetCursorPos")
	setCursorPos     = user32.NewProc("SetCursorPos")
	mouseEvent       = user32.NewProc("mouse_event")
	vkKeyScanA       = user32.NewProc("VkKeyScanA")
	keybdEvent       = user32.NewProc("keybd_event")
	getSystemMetrics = user32.NewProc("GetSystemMetrics")
)

type mouseCursor struct {
	x int
	y int
}

const (
	MOUSEEVENTF_MOVE        = 0x0001
	MOUSEEVENTF_LEFTDOWN    = 0x0002
	MOUSEEVENTF_LEFTUP      = 0x0004
	MOUSEEVENTF_LEFTCLICK   = MOUSEEVENTF_LEFTDOWN + MOUSEEVENTF_LEFTUP
	MOUSEEVENTF_RIGHTDOWN   = 0x0008
	MOUSEEVENTF_RIGHTUP     = 0x0010
	MOUSEEVENTF_RIGHTCLICK  = MOUSEEVENTF_RIGHTDOWN + MOUSEEVENTF_RIGHTUP
	MOUSEEVENTF_MIDDLEDOWN  = 0x0020
	MOUSEEVENTF_MIDDLEUP    = 0x0040
	MOUSEEVENTF_MIDDLECLICK = MOUSEEVENTF_MIDDLEDOWN + MOUSEEVENTF_MIDDLEUP

	MOUSEEVENTF_WHEEL = 0x0800

	KEYEVENTF_KEYDOWN = 0x0000
	KEYEVENTF_KEYUP   = 0x0002
)

func getKeyValue(key EventEntity) uintptr {
	if key >= Key0 && key <= Key9 {
		char := (key - Key0) + 48
		fmt.Println(string(char))
		keyCode, _, _ := vkKeyScanA.Call(uintptr(char))
		return keyCode
	} else if key >= KeyA && key <= KeyZ {
		char := (key - KeyA) + 97
		fmt.Println(string(char))
		keyCode, _, _ := vkKeyScanA.Call(uintptr(char))
		return keyCode
	} else {
		return func(key EventEntity) uintptr {
			return map[EventEntity]uintptr{
				KeyBackspace:     0x08, // VK_BACK
				KeyWinCmd:        0x5B, //VK_LWIN
				KeyTab:           0x09, // VK_TAB
				KeyEnter:         0x0d, // VK_RETURN
				KeyShift:         0x10, // VK_SHIFT
				KeyCtrl:          0x11, // VK_CONTROL
				KeyAlt:           0x12, // VK_MENU
				KeyPause:         0x13, // VK_PAUSE
				KeyCaps:          0x14, // VK_CAPITAL
				KeyEsc:           0x1b, // VK_ESCAPE
				KeySpace:         0x20, // VK_SPACE
				KeyPageUp:        0x21, // VK_PRIOR
				KeyPageDown:      0x22, // VK_NEXT
				KeyEnd:           0x23, // VK_END
				KeyHome:          0x24, // VK_HOME
				KeyArrowLeft:     0x25, // VK_LEFT
				KeyArrowUp:       0x26, // VK_UP
				KeyArrowRight:    0x27, // VK_RIGHT
				KeyArrowDown:     0x28, // VK_DOWN
				KeyPrintScreen:   0x2c, // VK_SNAPSHOT
				KeyInsert:        0x2d, // VK_INSERT
				KeyDelete:        0x2e, // VK_DELETE
				KeyLWinCmd:       0x5b, // VK_LWIN
				KeyRWinCmd:       0x5c, // VK_RWIN
				KeyNumPad0:       0x60, // VK_NUMPAD0
				KeyNumPad1:       0x61, // VK_NUMPAD1
				KeyNumPad2:       0x62, // VK_NUMPAD2
				KeyNumPad3:       0x63, // VK_NUMPAD3
				KeyNumPad4:       0x64, // VK_NUMPAD4
				KeyNumPad5:       0x65, // VK_NUMPAD5
				KeyNumPad6:       0x66, // VK_NUMPAD6
				KeyNumPad7:       0x67, // VK_NUMPAD7
				KeyNumPad8:       0x68, // VK_NUMPAD8
				KeyNumPad9:       0x69, // VK_NUMPAD9
				KeyNumPadMul:     0x6a, // VK_MULTIPLY  ??? Is this the numpad *?
				KeyNumPadAdd:     0x6b, // VK_ADD  ??? Is this the numpad +?
				KeyNumPadEnter:   0x6c, // VK_SEPARATOR  ??? Is this the numpad enter?
				KeyNumPadSub:     0x6d, // VK_SUBTRACT  ??? Is this the numpad -?
				KeyNumPadDecimal: 0x6e, // VK_DECIMAL
				KeyNumPadDiv:     0x6f, // VK_DIVIDE
				KeyF1:            0x70, // VK_F1
				KeyF2:            0x71, // VK_F2
				KeyF3:            0x72, // VK_F3
				KeyF4:            0x73, // VK_F4
				KeyF5:            0x74, // VK_F5
				KeyF6:            0x75, // VK_F6
				KeyF7:            0x76, // VK_F7
				KeyF8:            0x77, // VK_F8
				KeyF9:            0x78, // VK_F9
				KeyF10:           0x79, // VK_F10
				KeyF11:           0x7a, // VK_F11
				KeyF12:           0x7b, // VK_F12
				KeyNumLock:       0x90, // VK_NUMLOCK
				KeyScrollLock:    0x91, // VK_SCROLL
				KeyLShift:        0xa0, // VK_LSHIFT
				KeyRShift:        0xa1, // VK_RSHIFT
				KeyLCtrl:         0xa2, // VK_LCONTROL
				KeyRCtrl:         0xa3, // VK_RCONTROL
				KeyLAlt:          0xa4, // VK_LMENU
				KeyRAlt:          0xa5, // VK_RMENU
				KeyVolumeMute:    0xad, // VK_VOLUME_MUTE
				KeyVolumeDown:    0xae, // VK_VOLUME_DOWN
				KeyVolumeUp:      0xaf, // VK_VOLUME_UP
				KeyMediaNext:     0xb0, // VK_MEDIA_NEXT_TRACK
				KeyMediaPrevious: 0xb1, // VK_MEDIA_PREV_TRACK
				KeyMediaStop:     0xb2, // VK_MEDIA_STOP
				KeyMediaPause:    0xb3, // VK_MEDIA_PLAY_PAUSE
			}[key]
		}(key)
	}
}

func GetCursorPos() (int, int) {
	cursor := mouseCursor{}
	getCursorPos.Call(uintptr(unsafe.Pointer(&cursor)))
	return cursor.x, cursor.y
}

func GetScreenSize() (uintptr, uintptr) {
	width, _, _ := getSystemMetrics.Call(0)
	heigh, _, _ := getSystemMetrics.Call(1)
	return width, heigh
}

func SetCursorPos(x, y uintptr) uintptr {
	ret, _, _ := setCursorPos.Call(x, y)
	return ret
}

func MouseMove(x, y int) {
	mouseEvent.Call(MOUSEEVENTF_MOVE, uintptr(x), uintptr(y), 0)
}

func MouseScroll(lines int) {
	x, y := GetCursorPos()
	fmt.Println(x, y)
	mouseEvent.Call(MOUSEEVENTF_WHEEL, uintptr(x), uintptr(y), uintptr(lines))
}

func MouseRelease(button EventEntity) {
	x, y := GetCursorPos()
	var eventID uintptr
	switch button {
	case MouseLeftButton:
		eventID = MOUSEEVENTF_LEFTUP
	case MouseRightButton:
		eventID = MOUSEEVENTF_RIGHTUP
	case MouseMiddleButton:
		eventID = MOUSEEVENTF_MIDDLEUP
	}
	mouseEvent.Call(eventID, uintptr(x), uintptr(y), 0)
}

func MouseHold(button EventEntity) {
	x, y := GetCursorPos()
	var eventID uintptr
	switch button {
	case MouseLeftButton:
		eventID = MOUSEEVENTF_LEFTDOWN
	case MouseRightButton:
		eventID = MOUSEEVENTF_RIGHTDOWN
	case MouseMiddleButton:
		eventID = MOUSEEVENTF_MIDDLEDOWN
	}
	mouseEvent.Call(eventID, uintptr(x), uintptr(y), 0)
}

func MouseClick(button EventEntity) {
	x, y := GetCursorPos()
	var eventID uintptr
	switch button {
	case MouseLeftButton:
		eventID = MOUSEEVENTF_LEFTCLICK
	case MouseRightButton:
		eventID = MOUSEEVENTF_RIGHTCLICK
	case MouseMiddleButton:
		eventID = MOUSEEVENTF_MIDDLECLICK
	}
	mouseEvent.Call(eventID, uintptr(x), uintptr(y), 0)
}

func KeyHold(keyCode uintptr) uintptr {
	ret, _, _ := keybdEvent.Call(keyCode, 0, KEYEVENTF_KEYDOWN, 0)
	return ret
}

func KeyRelease(keyCode uintptr) uintptr {
	ret, _, _ := keybdEvent.Call(keyCode, 0, KEYEVENTF_KEYUP, 0)
	return ret
}

func init() {
	fmt.Print("Starting Up\n")
}
