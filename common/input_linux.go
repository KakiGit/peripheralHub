// +build linux

package common

import (
	"fmt"
)

type mouseCursor struct {
	x int
	y int
}

type Input struct {
	Display *Display
	com     chan InternalMsg
}

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

func getKeyboardEvent(event Event) int {
	return map[Event]int{
		ButtonDown: KeyPress,
		ButtonUp:   KeyRelease,
	}[event]
}

func getMouseEvent(event Event) int {
	return map[Event]int{
		ButtonDown: ButtonPress,
		ButtonUp:   ButtonRelease,
	}[event]
}

func getMouseButtonId(button EventEntity) uint {
	return map[EventEntity]uint{
		MouseLeftButton:   0x01,
		MouseMiddleButton: 0x02,
		MouseRightButton:  0x03,
		MouseOptButton1:   0x08,
		MouseOptButton2:   0x09,
	}[button]
}

func getKeyValue(display *Display, key EventEntity) uint {
	if key >= Key0 && key <= Key9 {
		char := (key - Key0) + 48
		keySym := XStringToKeysym(string(char))
		keyCode := XKeysymToKeycode(display, keySym)
		return keyCode
	} else if key >= KeyA && key <= KeyZ {
		char := (key - KeyA) + 97
		keySym := XStringToKeysym(string(char))
		keyCode := XKeysymToKeycode(display, keySym)
		return keyCode
	} else {
		return map[EventEntity]uint{
			KeyBackspace:     XKeysymToKeycode(display, XStringToKeysym(string("BackSpace"))),
			KeyWinCmd:        XKeysymToKeycode(display, XStringToKeysym(string("Super_L"))),
			KeyTab:           XKeysymToKeycode(display, XStringToKeysym(string("Tab"))),
			KeyEnter:         XKeysymToKeycode(display, XStringToKeysym(string("Return"))),
			KeyShift:         XKeysymToKeycode(display, XStringToKeysym(string("Shift_L"))),
			KeyCtrl:          XKeysymToKeycode(display, XStringToKeysym(string("Control_L"))),
			KeyAlt:           XKeysymToKeycode(display, XStringToKeysym(string("Alt_L"))),
			KeyPause:         XKeysymToKeycode(display, XStringToKeysym(string("Pause"))),
			KeyCaps:          XKeysymToKeycode(display, XStringToKeysym(string("Caps_Lock"))),
			KeyEsc:           XKeysymToKeycode(display, XStringToKeysym(string("Escape"))),
			KeySpace:         XKeysymToKeycode(display, XStringToKeysym(string("space"))),
			KeyPageUp:        XKeysymToKeycode(display, XStringToKeysym(string("Prior"))),
			KeyPageDown:      XKeysymToKeycode(display, XStringToKeysym(string("Next"))),
			KeyEnd:           XKeysymToKeycode(display, XStringToKeysym(string("End"))),
			KeyHome:          XKeysymToKeycode(display, XStringToKeysym(string("Home"))),
			KeyArrowLeft:     XKeysymToKeycode(display, XStringToKeysym(string("Left"))),
			KeyArrowUp:       XKeysymToKeycode(display, XStringToKeysym(string("Up"))),
			KeyArrowRight:    XKeysymToKeycode(display, XStringToKeysym(string("Right"))),
			KeyArrowDown:     XKeysymToKeycode(display, XStringToKeysym(string("Down"))),
			KeyPrintScreen:   XKeysymToKeycode(display, XStringToKeysym(string("Print"))),
			KeyInsert:        XKeysymToKeycode(display, XStringToKeysym(string("Insert"))),
			KeyDelete:        XKeysymToKeycode(display, XStringToKeysym(string("Delete"))),
			KeyLWinCmd:       XKeysymToKeycode(display, XStringToKeysym(string("Super_L"))),
			KeyRWinCmd:       XKeysymToKeycode(display, XStringToKeysym(string("Super_R"))),
			KeyNumPad0:       XKeysymToKeycode(display, XStringToKeysym(string("KP_Insert"))),
			KeyNumPad1:       XKeysymToKeycode(display, XStringToKeysym(string("KP_End"))),
			KeyNumPad2:       XKeysymToKeycode(display, XStringToKeysym(string("KP_Down"))),
			KeyNumPad3:       XKeysymToKeycode(display, XStringToKeysym(string("KP_Next"))),
			KeyNumPad4:       XKeysymToKeycode(display, XStringToKeysym(string("KP_Left"))),
			KeyNumPad5:       XKeysymToKeycode(display, XStringToKeysym(string("KP_Begin"))),
			KeyNumPad6:       XKeysymToKeycode(display, XStringToKeysym(string("KP_Right"))),
			KeyNumPad7:       XKeysymToKeycode(display, XStringToKeysym(string("KP_Home"))),
			KeyNumPad8:       XKeysymToKeycode(display, XStringToKeysym(string("KP_Up"))),
			KeyNumPad9:       XKeysymToKeycode(display, XStringToKeysym(string("KP_Prior"))),
			KeyNumPadMul:     XKeysymToKeycode(display, XStringToKeysym(string("KP_Multiply"))),
			KeyNumPadAdd:     XKeysymToKeycode(display, XStringToKeysym(string("KP_Add"))),
			KeyNumPadEnter:   XKeysymToKeycode(display, XStringToKeysym(string("KP_Enter"))),
			KeyNumPadSub:     XKeysymToKeycode(display, XStringToKeysym(string("KP_Subtract"))),
			KeyNumPadDecimal: XKeysymToKeycode(display, XStringToKeysym(string("KP_Delete"))),
			KeyNumPadDiv:     XKeysymToKeycode(display, XStringToKeysym(string("KP_Divide"))),
			KeyF1:            XKeysymToKeycode(display, XStringToKeysym(string("F1"))),
			KeyF2:            XKeysymToKeycode(display, XStringToKeysym(string("F2"))),
			KeyF3:            XKeysymToKeycode(display, XStringToKeysym(string("F3"))),
			KeyF4:            XKeysymToKeycode(display, XStringToKeysym(string("F4"))),
			KeyF5:            XKeysymToKeycode(display, XStringToKeysym(string("F5"))),
			KeyF6:            XKeysymToKeycode(display, XStringToKeysym(string("F6"))),
			KeyF7:            XKeysymToKeycode(display, XStringToKeysym(string("F7"))),
			KeyF8:            XKeysymToKeycode(display, XStringToKeysym(string("F8"))),
			KeyF9:            XKeysymToKeycode(display, XStringToKeysym(string("F9"))),
			KeyF10:           XKeysymToKeycode(display, XStringToKeysym(string("F10"))),
			KeyF11:           XKeysymToKeycode(display, XStringToKeysym(string("F11"))),
			KeyF12:           XKeysymToKeycode(display, XStringToKeysym(string("F12"))),
			KeyNumLock:       XKeysymToKeycode(display, XStringToKeysym(string("Num_Lock"))),
			KeyScrollLock:    XKeysymToKeycode(display, XStringToKeysym(string("Scroll_Lock"))),
			KeyLShift:        XKeysymToKeycode(display, XStringToKeysym(string("Shift_L"))),
			KeyRShift:        XKeysymToKeycode(display, XStringToKeysym(string("Shift_R"))),
			KeyLCtrl:         XKeysymToKeycode(display, XStringToKeysym(string("Control_L"))),
			KeyRCtrl:         XKeysymToKeycode(display, XStringToKeysym(string("Control_R"))),
			KeyLAlt:          XKeysymToKeycode(display, XStringToKeysym(string("Alt_L"))),
			KeyRAlt:          XKeysymToKeycode(display, XStringToKeysym(string("Alt_R"))),
			KeyTilde:         XKeysymToKeycode(display, XStringToKeysym(string("grave"))),
			// KeyVolumeMute:    XKeysymToKeycode(display, XStringToKeysym(string("backspace"))),
			// KeyVolumeDown:    XKeysymToKeycode(display, XStringToKeysym(string("backspace"))),
			// KeyVolumeUp:      XKeysymToKeycode(display, XStringToKeysym(string("backspace"))),
			// KeyMediaNext:     XKeysymToKeycode(display, XStringToKeysym(string("backspace"))),
			// KeyMediaPrevious: XKeysymToKeycode(display, XStringToKeysym(string("backspace"))),
			// KeyMediaStop:     XKeysymToKeycode(display, XStringToKeysym(string("backspace"))),
			// KeyMediaPause:    XKeysymToKeycode(display, XStringToKeysym(string("backspace"))),
		}[key]
	}
}

func (input *Input) GetCursorPos() (int, int) {
	var event XEvent
	XNextEvent(input.Display, &event)
	x, y := GetCursorPosition(&event)
	return x, y
}

func (input *Input) SetCursorPos(x, y int) uintptr {
	XTestFakeMotionEvent(input.Display, 0, x, y, 0)
	return 1
}

func (input *Input) MouseMove(x, y int) {
	XTestFakeRelativeMotionEvent(input.Display, x, y, 0)
}

func (input *Input) MouseScroll(lines int) {
	var keycode uint
	if lines > 0 {
		keycode = 4
	} else {
		keycode = 5
	}
	for i := 0; i < lines; i++ {
		XTestFakeButtonEvent(input.Display, keycode, KeyPress, 0)
		XTestFakeButtonEvent(input.Display, keycode, KeyRelease, 0)
	}
}

func (input *Input) MouseButtonAction(button EventEntity, event Event) {
	eventID := getKeyboardEvent(event)
	fmt.Println(eventID)
	keycode := getMouseButtonId(button)
	fmt.Println(keycode)
	XTestFakeButtonEvent(input.Display, keycode, eventID, 0)
}

func (input *Input) KeyboardButtonAction(button EventEntity, event Event) {
	FakeKeyEvent(input.Display, getKeyValue(input.Display, button), getKeyboardEvent(event), 0)
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
	input.Display = XOpenDisplay()
	com := make(chan InternalMsg, 50)
	go func(com chan InternalMsg) {
		XTestGrabControl(input.Display, True)
		defer XTestGrabControl(input.Display, False)
		for {
			internalMsg := <-com
			event := internalMsg.Event
			eventEntity := internalMsg.EventEntity
			switch event {
			case ButtonDown:
				input.ButtonAction(eventEntity, event)
			case ButtonUp:
				input.ButtonAction(eventEntity, event)
			case MouseRelativeMove:
				input.MouseMove(internalMsg.ExtraInfo[0], internalMsg.ExtraInfo[1])
			case MouseWheelScrollUp:
				input.MouseScroll(1)
			case MouseWheelScrollDown:
				input.MouseScroll(-1)
			}
		}
	}(com)
}

func (input *Input) InputFromClient(message Message) {

	internalMsg := InternalMsg{
		EventEntity: message.EventEntity,
		Event:       message.Event,
		ExtraInfo:   message.ExtraInfo,
	}
	input.com <- internalMsg

}
