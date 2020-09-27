// +build windows

package common

import (
	"testing"
	"time"
)

func TestWinInput(t *testing.T) {
	input := Input{}
	input.SetCursorPos(2600, 900)
}

func TestKeyPress(t *testing.T) {
	input := Input{}
	input.KeyboardButtonAction(KeyShift, ButtonDown)
	input.KeyboardButtonAction(KeyA, ButtonDown)
	input.KeyboardButtonAction(KeyA, ButtonUp)
	input.KeyboardButtonAction(KeyShift, ButtonUp)
}

func TestScroll(t *testing.T) {
	input := Input{}
	input.MouseScroll(100)
}

func TestMouseClick(t *testing.T) {
	input := Input{}
	input.MouseButtonAction(MouseLeftButton, ButtonDown)
	time.Sleep(time.Second * 2)
	input.MouseButtonAction(MouseLeftButton, ButtonUp)
}

func TestMouseMove(t *testing.T) {
	input := Input{}
	input.MouseMove(100, 100)
}
