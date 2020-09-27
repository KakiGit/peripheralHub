// +build windows

package common

import (
	"fmt"
	"syscall"
	"unsafe"
)

type Input struct {
}

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

func getKeyboardEvent(event Event) uintptr {
	return map[Event]uintptr{
		ButtonDown: 0x0000,
		ButtonUp:   0x0002,
	}[event]
}

func getMouseEventID(button EventEntity, event Event) uintptr {
	return map[EventEntity]map[Event]uintptr{
		MouseLeftButton: map[Event]uintptr{
			ButtonDown:  0x0002,
			ButtonUp:    0x0004,
			ButtonClick: 0x0006,
		},
		MouseRightButton: map[Event]uintptr{
			ButtonDown:  0x0008,
			ButtonUp:    0x0010,
			ButtonClick: 0x0018,
		},
		MouseMiddleButton: map[Event]uintptr{
			ButtonDown:  0x0020,
			ButtonUp:    0x0040,
			ButtonClick: 0x0060,
		},
		MouseWheel: map[Event]uintptr{
			MouseWheelScroll: 0x0800,
		},
		MouseCursor: map[Event]uintptr{
			MouseRelativeMove: 0x0001,
		},
	}[button][event]
}

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
	}
}

func (input *Input) GetCursorPos() (int, int) {
	cursor := mouseCursor{}
	getCursorPos.Call(uintptr(unsafe.Pointer(&cursor)))
	return cursor.x, cursor.y
}

func (input *Input) GetScreenSize() (uintptr, uintptr) {
	width, _, _ := getSystemMetrics.Call(0)
	heigh, _, _ := getSystemMetrics.Call(1)
	return width, heigh
}

func (input *Input) SetCursorPos(x, y uintptr) {
	setCursorPos.Call(x, y)
}

func (input *Input) MouseMove(x, y int) {
	mouseEvent.Call(getMouseEventID(MouseCursor, MouseRelativeMove), uintptr(x), uintptr(y), 0)
}

func (input *Input) MouseScroll(lines int) {
	x, y := input.GetCursorPos()
	fmt.Println(x, y)
	mouseEvent.Call(getMouseEventID(MouseWheel, MouseWheelScroll), uintptr(x), uintptr(y), uintptr(lines))
}

func (input *Input) MouseButtonAction(button EventEntity, event Event) {
	x, y := input.GetCursorPos()
	eventID := getMouseEventID(button, event)
	mouseEvent.Call(eventID, uintptr(x), uintptr(y), 0)
}

func (input *Input) KeyboardButtonAction(button EventEntity, event Event) {
	keybdEvent.Call(getKeyValue(button), 0, getKeyboardEvent(event), 0)
}

func (input *Input) ButtonAction(button EventEntity, event Event) {
	if button == MouseLeftButton || button == MouseMiddleButton || button == MouseRightButton {
		input.MouseButtonAction(button, event)
	} else {
		input.KeyboardButtonAction(button, event)
	}
}

func (input *Input) Init() {
	fmt.Print("Starting Up\n")
}

func (input *Input) InputFromClient(message Message) {
	event := message.Event
	eventEntity := message.EventEntity
	switch event {
	case ButtonDown:
		input.ButtonAction(eventEntity, event)
	case ButtonUp:
		input.ButtonAction(eventEntity, event)
	case MouseRelativeMove:
		input.MouseMove(message.ExtraInfo[0], message.ExtraInfo[1])
	case MouseWheelScrollUp:
		input.MouseScroll(10)
	case MouseWheelScrollDown:
		input.MouseScroll(-10)
	}
}
