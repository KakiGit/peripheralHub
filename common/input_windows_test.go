// +build windows
package common

import (
	"fmt"
	"testing"
	"time"
)

func TestWinInput(t *testing.T) {
	// defer syscall.FreeLibrary(user32)

	// fmt.Printf("Return: %d\n", MessageBox("Done Title", "This test is Done.", MB_YESNOCANCEL))
	fmt.Printf("Return: %d\n", SetCursorPos(2600, 900))
}

func TestKeyPress(t *testing.T) {
	// defer syscall.FreeLibrary(user32)
	keyCode := getKeyValue(KeyA)
	keyCode2 := getKeyValue(KeyLShift)
	res := KeyHold(keyCode2)
	res = KeyHold(keyCode)
	fmt.Println(res)
	res = KeyRelease(keyCode)
	res = KeyRelease(keyCode2)
	fmt.Println(res)
}

func TestScroll(t *testing.T) {
	MouseScroll(-10)
}

func TestMouseClick(t *testing.T) {
	MouseButtonAction(MouseLeftButton, ButtonDown)
	time.Sleep(time.Second * 2)
	MouseButtonAction(MouseLeftButton, ButtonUp)
}

func TestMouseMove(t *testing.T) {
	MouseMove(100, 100)
}
