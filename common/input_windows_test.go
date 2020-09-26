// +build windows

package common

import (
	"testing"
	"time"
)

func TestWinInput(t *testing.T) {

	SetCursorPos(2600, 900)
}

func TestKeyPress(t *testing.T) {
	KeyHold(KeyShift)
	KeyHold(KeyA)
	KeyRelease(KeyA)
	KeyRelease(KeyShift)
}

func TestScroll(t *testing.T) {
	MouseScroll(100)
}

func TestMouseClick(t *testing.T) {
	MouseButtonAction(MouseLeftButton, ButtonDown)
	time.Sleep(time.Second * 2)
	MouseButtonAction(MouseLeftButton, ButtonUp)
}

func TestMouseMove(t *testing.T) {
	MouseMove(100, 100)
}
