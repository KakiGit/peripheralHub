// +build linux

package common

import (
	"testing"
	"time"
)

func TestKeyPress(t *testing.T) {
	input := Input{}
	display := XOpenDisplay()
	input.Display = display
	XTestGrabControl(display, True)
	input.KeyHold(KeyShift)
	input.KeyHold(KeyA)
	input.KeyRelease(KeyA)
	input.KeyRelease(KeyShift)
	XTestGrabControl(display, False)
}

func TestScroll(t *testing.T) {
	input := Input{}
	display := XOpenDisplay()
	input.Display = display
	XTestGrabControl(display, True)
	input.MouseScroll(-10)
	time.Sleep(1 * time.Second)
	input.MouseScroll(10)
	XTestGrabControl(display, False)
}

func TestMouseClick(t *testing.T) {
	input := Input{}
	display := XOpenDisplay()
	input.Display = display
	XTestGrabControl(display, True)
	input.MouseButtonAction(MouseLeftButton, ButtonDown)
	time.Sleep(time.Second * 2)
	input.MouseButtonAction(MouseLeftButton, ButtonUp)
	XTestGrabControl(display, False)
}

func TestMouseMove(t *testing.T) {
	input := Input{}
	display := XOpenDisplay()
	input.Display = display
	XTestGrabControl(display, True)
	input.MouseMove(100, 100)
	XTestGrabControl(display, False)
}
