// +build linux

package common

import (
	"fmt"
)

type Sender struct {
	Address         string
	Devices         []Device
	Receivers       []*Receiver
	CurrentReceiver *Receiver
	Secret          Secret
}

type InternalMsg struct {
	EventEntity EventEntity
	Event       Event
}

func getEventFromKeycode(keycode uint, xeventType int) Event {
	if keycode == 4 || keycode == 5 {
		return MouseWheelScroll
	}
	return map[int]Event{
		ButtonPress:   ButtonDown,
		ButtonRelease: ButtonUp,
	}[xeventType]
}

func getEventEntityFromMouseKeycode(keycode uint) EventEntity {
	return map[uint]EventEntity{
		1: MouseLeftButton,
		2: MouseMiddleButton,
		3: MouseRightButton,
		4: MouseWheel,
		5: MouseWheel,
	}[keycode]
}

func getEventEntityFromKeyString(keyString string) EventEntity {
	return map[string]EventEntity{
		"BackSpace":   KeyBackspace,
		"Tab":         KeyTab,
		"Return":      KeyEnter,
		"Pause":       KeyPause,
		"Caps_Lock":   KeyCaps,
		"Escape":      KeyEsc,
		"space":       KeySpace,
		"Prior":       KeyPageUp,
		"Next":        KeyPageDown,
		"End":         KeyEnd,
		"Home":        KeyHome,
		"Left":        KeyArrowLeft,
		"Up":          KeyArrowUp,
		"Right":       KeyArrowRight,
		"Down":        KeyArrowDown,
		"Print":       KeyPrintScreen,
		"Insert":      KeyInsert,
		"Delete":      KeyDelete,
		"Super_L":     KeyLWinCmd,
		"Super_R":     KeyRWinCmd,
		"KP_Insert":   KeyNumPad0,
		"KP_End":      KeyNumPad1,
		"KP_Down":     KeyNumPad2,
		"KP_Next":     KeyNumPad3,
		"KP_Left":     KeyNumPad4,
		"KP_Begin":    KeyNumPad5,
		"KP_Right":    KeyNumPad6,
		"KP_Home":     KeyNumPad7,
		"KP_Up":       KeyNumPad8,
		"KP_Prior":    KeyNumPad9,
		"KP_Multiply": KeyNumPadMul,
		"KP_Add":      KeyNumPadAdd,
		"KP_Enter":    KeyNumPadEnter,
		"KP_Subtract": KeyNumPadSub,
		"KP_Delete":   KeyNumPadDecimal,
		"KP_Divide":   KeyNumPadDiv,
		"F1":          KeyF1,
		"F2":          KeyF2,
		"F3":          KeyF3,
		"F4":          KeyF4,
		"F5":          KeyF5,
		"F6":          KeyF6,
		"F7":          KeyF7,
		"F8":          KeyF8,
		"F9":          KeyF9,
		"F10":         KeyF10,
		"F11":         KeyF11,
		"F12":         KeyF12,
		"Num_Lock":    KeyNumLock,
		"Scroll_Lock": KeyScrollLock,
		"Shift_L":     KeyLShift,
		"Shift_R":     KeyRShift,
		"Control_L":   KeyLCtrl,
		"Control_R":   KeyRCtrl,
		"Alt_L":       KeyLAlt,
		"Alt_R":       KeyRAlt,
		"0":           Key0,
		"1":           Key1,
		"2":           Key2,
		"3":           Key3,
		"4":           Key4,
		"5":           Key5,
		"6":           Key6,
		"7":           Key7,
		"8":           Key8,
		"9":           Key9,
		"a":           KeyA,
		"b":           KeyB,
		"c":           KeyC,
		"d":           KeyD,
		"e":           KeyE,
		"f":           KeyF,
		"g":           KeyG,
		"h":           KeyH,
		"i":           KeyI,
		"j":           KeyJ,
		"k":           KeyK,
		"l":           KeyL,
		"m":           KeyM,
		"n":           KeyN,
		"o":           KeyO,
		"p":           KeyP,
		"q":           KeyQ,
		"r":           KeyR,
		"s":           KeyS,
		"t":           KeyT,
		"u":           KeyU,
		"v":           KeyV,
		"w":           KeyW,
		"x":           KeyX,
		"y":           KeyY,
		"z":           KeyZ,
	}[keyString]

}

func ListenXEvent(com chan InternalMsg) {
	display := XOpenDisplay()
	root := XDefaultRootWindow(display)
	XGrabPointer(display, root)
	XGrabKeyboard(display, root)
	var event XEvent
	for {
		XNextEvent(display, &event)
		eventType := GetXEventType(&event)
		switch eventType {
		case KeyPress:
			keycode := GetKeyCode(&event)
			keysym := XKeycodeToKeysym(display, keycode)
			keystring := XKeysymToString(keysym)
			fmt.Printf("code to string -- keycode: %v, keysym: %v, string: %v\n", keycode, keysym, keystring)
			fmt.Printf("KeyPress %d, 0x%x, %v\n", keycode, keycode, eventType)
			com <- InternalMsg{
				EventEntity: getEventEntityFromKeyString(keystring),
				Event:       ButtonDown,
			}
		case KeyRelease:
			keycode := GetKeyCode(&event)
			keysym := XKeycodeToKeysym(display, keycode)
			keystring := XKeysymToString(keysym)
			fmt.Printf("code to string -- keycode: %v, keysym: %v, string: %v\n", keycode, keysym, keystring)
			fmt.Printf("KeyRelease %d, 0x%x, %v\n", keycode, keycode, eventType)
			com <- InternalMsg{
				EventEntity: getEventEntityFromKeyString(keystring),
				Event:       ButtonUp,
			}
		case ButtonPress:
			keycode := GetKeyCode(&event)
			keysym := XKeycodeToKeysym(display, keycode)
			keystring := XKeysymToString(keysym)
			fmt.Printf("code to string -- keycode: %v, keysym: %v, string: %v\n", keycode, keysym, keystring)
			fmt.Printf("ButtonPress %d, 0x%x, %v\n", keycode, keycode, eventType)
			com <- InternalMsg{
				EventEntity: getEventEntityFromMouseKeycode(keycode),
				Event:       getEventFromKeycode(keycode, ButtonPress),
			}
		case ButtonRelease:
			keycode := GetKeyCode(&event)
			keysym := XKeycodeToKeysym(display, keycode)
			keystring := XKeysymToString(keysym)
			fmt.Printf("code to string -- keycode: %v, keysym: %v, string: %v\n", keycode, keysym, keystring)
			fmt.Printf("ButtonRelease %d, 0x%x, %v\n", keycode, keycode, eventType)
			com <- InternalMsg{
				EventEntity: getEventEntityFromKeyString(keystring),
				Event:       getEventFromKeycode(keycode, ButtonRelease),
			}
		case MotionNotify:
			x, y := GetCursorPosition(&event)
			fmt.Printf("MotionNotify x:%d, 0x%x y:%d, 0x%x, %v\n", x, x, y, y, eventType)
			com <- InternalMsg{
				EventEntity: MouseCursor,
				Event:       MouseRelativeMove,
			}
		}
	}
}
