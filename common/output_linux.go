// +build linux

package common

import (
	"fmt"

	evdev "github.com/gvalkov/golang-evdev"
)

type Output struct {
	display *Display
	root    Window
	Com     chan InternalMsg
}

// func getEventEntityAndEvent(value int, eventType int) (EventEntity, Event) {
// 	switch
// }

func getEventEntityFromMouseKeycode(keycode uint) (EventEntity, bool) {
	eventEntities := map[uint]EventEntity{
		evdev.BTN_LEFT:   MouseLeftButton,
		evdev.BTN_MIDDLE: MouseMiddleButton,
		evdev.BTN_RIGHT:  MouseRightButton,
	}
	eventEntity, ok := eventEntities[keycode]
	return eventEntity, ok
}

func getEventEntityFromKeyString(keyString string) EventEntity {
	return map[string]EventEntity{
		"BackSpace":    KeyBackspace,
		"Tab":          KeyTab,
		"Return":       KeyEnter,
		"Pause":        KeyPause,
		"Caps_Lock":    KeyCaps,
		"Escape":       KeyEsc,
		"space":        KeySpace,
		"Prior":        KeyPageUp,
		"Next":         KeyPageDown,
		"End":          KeyEnd,
		"Home":         KeyHome,
		"Left":         KeyArrowLeft,
		"Up":           KeyArrowUp,
		"Right":        KeyArrowRight,
		"Down":         KeyArrowDown,
		"Print":        KeyPrintScreen,
		"Insert":       KeyInsert,
		"Delete":       KeyDelete,
		"Super_L":      KeyLWinCmd,
		"Super_R":      KeyRWinCmd,
		"KP_Insert":    KeyNumPad0,
		"KP_End":       KeyNumPad1,
		"KP_Down":      KeyNumPad2,
		"KP_Next":      KeyNumPad3,
		"KP_Left":      KeyNumPad4,
		"KP_Begin":     KeyNumPad5,
		"KP_Right":     KeyNumPad6,
		"KP_Home":      KeyNumPad7,
		"KP_Up":        KeyNumPad8,
		"KP_Prior":     KeyNumPad9,
		"KP_Multiply":  KeyNumPadMul,
		"KP_Add":       KeyNumPadAdd,
		"KP_Enter":     KeyNumPadEnter,
		"KP_Subtract":  KeyNumPadSub,
		"KP_Delete":    KeyNumPadDecimal,
		"KP_Divide":    KeyNumPadDiv,
		"grave":        KeyTilde,
		"F1":           KeyF1,
		"F2":           KeyF2,
		"F3":           KeyF3,
		"F4":           KeyF4,
		"F5":           KeyF5,
		"F6":           KeyF6,
		"F7":           KeyF7,
		"F8":           KeyF8,
		"F9":           KeyF9,
		"F10":          KeyF10,
		"F11":          KeyF11,
		"F12":          KeyF12,
		"Num_Lock":     KeyNumLock,
		"Scroll_Lock":  KeyScrollLock,
		"Shift_L":      KeyLShift,
		"Shift_R":      KeyRShift,
		"Control_L":    KeyLCtrl,
		"Control_R":    KeyRCtrl,
		"Alt_L":        KeyLAlt,
		"Alt_R":        KeyRAlt,
		"0":            Key0,
		"1":            Key1,
		"2":            Key2,
		"3":            Key3,
		"4":            Key4,
		"5":            Key5,
		"6":            Key6,
		"7":            Key7,
		"8":            Key8,
		"9":            Key9,
		"a":            KeyA,
		"b":            KeyB,
		"c":            KeyC,
		"d":            KeyD,
		"e":            KeyE,
		"f":            KeyF,
		"g":            KeyG,
		"h":            KeyH,
		"i":            KeyI,
		"j":            KeyJ,
		"k":            KeyK,
		"l":            KeyL,
		"m":            KeyM,
		"n":            KeyN,
		"o":            KeyO,
		"p":            KeyP,
		"q":            KeyQ,
		"r":            KeyR,
		"s":            KeyS,
		"t":            KeyT,
		"u":            KeyU,
		"v":            KeyV,
		"w":            KeyW,
		"x":            KeyX,
		"y":            KeyY,
		"z":            KeyZ,
		"minus":        KeyMinus,
		"equal":        KeyEqual,
		"bracketleft":  KeyLeftBracket,
		"backslash":    KeyBackSlash,
		"bracketright": KeyRightBracket,
		"apostrophe":   KeyApostrophe,
		"semicolon":    KeySemicolon,
		"comma":        KeyComma,
		"period":       KeyPeriod,
		"slash":        KeySlash,
	}[keyString]

}

func (output *Output) Init() {
	output.display = XOpenDisplay()
	output.root = XDefaultRootWindow(output.display)
	output.Com = make(chan InternalMsg, 50)
}

func (output *Output) OutputToServer() {
	go grabMouseEvent(output.Com)
	var event XEvent
	xReset := 250
	yReset := 200
	XGrabKeyboard(output.display, output.root)
	XTestGrabControl(output.display, True)
	XTestFakeMotionEvent(output.display, 0, xReset, yReset, 0)
	XNextEvent(output.display, &event)
	XTestGrabControl(output.display, False)
	for {
		XNextEvent(output.display, &event)
		eventType := GetXEventType(&event)
		switch eventType {
		case KeyPress:
			keycode := GetKeyCode(&event)
			keysym := XKeycodeToKeysym(output.display, keycode)
			keystring := XKeysymToString(keysym)

			output.Com <- InternalMsg{
				EventEntity: getEventEntityFromKeyString(keystring),
				Event:       ButtonDown,
			}
		case KeyRelease:
			keycode := GetKeyCode(&event)
			keysym := XKeycodeToKeysym(output.display, keycode)
			keystring := XKeysymToString(keysym)
			fmt.Printf("%v %v %v", keycode, keysym, keystring)
			output.Com <- InternalMsg{
				EventEntity: getEventEntityFromKeyString(keystring),
				Event:       ButtonUp,
			}
		}
	}
}

func grabMouseEvent(com chan InternalMsg) {
	devices, _ := evdev.ListInputDevices()

	for _, dev := range devices {
		fmt.Printf("%s %s %s\n", dev.Fn, dev.Name, dev.Phys)
	}
	device, _ := evdev.Open(devices[0].Fn)
	device.SetRepeatRate(0, 0)
	fmt.Println(device)
	device.Grab()
	defer device.Release()
	for {
		r, _ := device.Read()
		internalMsg := InternalMsg{}
		fmt.Println(r)
		for _, ev := range r {
			switch ev.Type {
			case evdev.EV_REL:
				if ev.Code == 8 {
					internalMsg.EventEntity = MouseWheel
					switch line := r[1].Value; {
					case line > 0:
						internalMsg.Event = MouseWheelScrollUp
					case line < 0:
						internalMsg.Event = MouseWheelScrollDown
					}
				} else if ev.Code == 0 || ev.Code == 1 {
					internalMsg.EventEntity = MouseCursor
					internalMsg.Event = MouseRelativeMove
					internalMsg.ExtraInfo[ev.Code] = int(ev.Value)
				}
			case evdev.EV_KEY:
				eventEntity, ok := getEventEntityFromMouseKeycode(uint(ev.Code))
				if !ok {
					continue
				}
				internalMsg.EventEntity = eventEntity
				if ev.Value == 1 {
					internalMsg.Event = ButtonDown
				} else {
					internalMsg.Event = ButtonUp
				}
			}
		}
		com <- internalMsg
	}
}
