// +build linux

package common

import (
	"fmt"
	"testing"
	"time"
)

func TestXlib(t *testing.T) {
	display := XOpenDisplay()
	keys := []string{"a", "b", "c", "d", "e"}

	XTestGrabControl(display, True)
	for _, key := range keys {
		keysym := XStringToKeysym(key)
		keycode := XKeysymToKeycode(display, keysym)
		fmt.Printf("String to code -- keycode: %v, keysym: %v, string: %v\n", keycode, keysym, key)
		FakeKeyEvent(display, keycode, KeyPress, 0)
		FakeKeyEvent(display, keycode, KeyRelease, 0)
		keysym = XKeycodeToKeysym(display, keycode)
		keystring := XKeysymToString(keysym)
		fmt.Printf("code to string -- keycode: %v, keysym: %v, string: %v\n", keycode, keysym, keystring)
	}
	XTestFakeButtonEvent(display, 0x1, KeyPress, 0)
	time.Sleep(2 * time.Second)
	XTestFakeButtonEvent(display, 0x1, KeyRelease, 0)

	XTestFakeMotionEvent(display, 0, 100, 100, 0)
	time.Sleep(2 * time.Second)

	XTestFakeRelativeMotionEvent(display, 100, 100, 0)
	XTestGrabControl(display, False)
}
