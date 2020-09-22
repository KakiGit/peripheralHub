package common

import (
	"fmt"
	"syscall"
)

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

var (
	user32, _       = syscall.LoadLibrary("user32.dll")
	setCursorPos, _ = syscall.GetProcAddress(user32, "SetCursorPos")
)

const (
	MOUSEEVENTF_MOVE        = 0x0001
	MOUSEEVENTF_LEFTDOWN    = 0x0002
	MOUSEEVENTF_LEFTUP      = 0x0004
	MOUSEEVENTF_LEFTCLICK   = MOUSEEVENTF_LEFTDOWN + MOUSEEVENTF_LEFTUP
	MOUSEEVENTF_RIGHTDOWN   = 0x0008
	MOUSEEVENTF_RIGHTUP     = 0x0010
	MOUSEEVENTF_RIGHTCLICK  = MOUSEEVENTF_RIGHTDOWN + MOUSEEVENTF_RIGHTUP
	MOUSEEVENTF_MIDDLEDOWN  = 0x0020
	MOUSEEVENTF_MIDDLEUP    = 0x0040
	MOUSEEVENTF_MIDDLECLICK = MOUSEEVENTF_MIDDLEDOWN + MOUSEEVENTF_MIDDLEUP

	MOUSEEVENTF_ABSOLUTE = 0x8000
	MOUSEEVENTF_WHEEL    = 0x0800
	MOUSEEVENTF_HWHEEL   = 0x01000

	KEYEVENTF_KEYDOWN = 0x0000 // Technically this constant doesn't exist in the MS documentation. It's the lack of KEYEVENTF_KEYUP that means pressing the key down.
	KEYEVENTF_KEYUP   = 0x0002

	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1
)

func SetCursorPos(x, y uintptr) (result int) {
	var nargs uintptr = 2
	ret, _, callErr := syscall.Syscall9(uintptr(setCursorPos),
		nargs,
		x,
		y,
		0,
		0,
		0,
		0,
		0,
		0,
		0)
	if callErr != 0 {
		abort("Call SetCursorPos", callErr)
	}
	result = int(ret)
	fmt.Println(result)
	return
}

func init() {
	fmt.Print("Starting Up\n")
}
