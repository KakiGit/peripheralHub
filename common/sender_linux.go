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

func ListenXEvent() {
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
		case KeyRelease:
			keycode := GetKeyCode(&event)
			keysym := XKeycodeToKeysym(display, keycode)
			keystring := XKeysymToString(keysym)
			fmt.Printf("code to string -- keycode: %v, keysym: %v, string: %v\n", keycode, keysym, keystring)
			fmt.Printf("KeyRelease %d, 0x%x, %v\n", keycode, keycode, eventType)
		case ButtonPress:
			keycode := GetKeyCode(&event)
			keysym := XKeycodeToKeysym(display, keycode)
			keystring := XKeysymToString(keysym)
			fmt.Printf("code to string -- keycode: %v, keysym: %v, string: %v\n", keycode, keysym, keystring)
			fmt.Printf("ButtonPress %d, 0x%x, %v\n", keycode, keycode, eventType)
		case ButtonRelease:
			keycode := GetKeyCode(&event)
			keysym := XKeycodeToKeysym(display, keycode)
			keystring := XKeysymToString(keysym)
			fmt.Printf("code to string -- keycode: %v, keysym: %v, string: %v\n", keycode, keysym, keystring)
			fmt.Printf("ButtonRelease %d, 0x%x, %v\n", keycode, keycode, eventType)
		case MotionNotify:
			x, y := GetCursorPosition(&event)
			fmt.Printf("MotionNotify x:%d, 0x%x y:%d, 0x%x, %v\n", x, x, y, y, eventType)
		}
	}
}
