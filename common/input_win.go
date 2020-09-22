package common

import (
	"fmt"
	"reflect"
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

func CallFunc(funcPointer uintptr, rawArgs ...uintptr) (r1, r2 uintptr, err syscall.Errno) {
	nargs := len(rawArgs)
	var args []reflect.Value
	args = append(args, reflect.ValueOf(funcPointer))
	args = append(args, reflect.ValueOf(uintptr(nargs)))
	for _, x := range rawArgs {
		args = append(args, reflect.ValueOf(x))
	}
	var sysc reflect.Value
	fillCount := 0
	switch {
	case nargs <= 3:
		sysc = reflect.ValueOf(syscall.Syscall)
		fillCount = 3 - nargs
	case nargs <= 6:
		sysc = reflect.ValueOf(syscall.Syscall6)
		fillCount = 6 - nargs
	case nargs <= 9:
		sysc = reflect.ValueOf(syscall.Syscall9)
		fillCount = 9 - nargs
	case nargs <= 12:
		sysc = reflect.ValueOf(syscall.Syscall12)
		fillCount = 12 - nargs
	case nargs <= 15:
		sysc = reflect.ValueOf(syscall.Syscall15)
		fillCount = 15 - nargs
	case nargs <= 18:
		sysc = reflect.ValueOf(syscall.Syscall18)
		fillCount = 18 - nargs
	}
	for fillCount > 0 {
		args = append(args, reflect.ValueOf(uintptr(0)))
		fillCount--
	}
	result := sysc.Call(args)
	r1 = result[0].Interface().(uintptr)
	r2 = result[1].Interface().(uintptr)
	err = result[2].Interface().(syscall.Errno)
	return
}

func SetCursorPos(x, y uintptr) (result int) {
	ret, _, callErr := CallFunc(setCursorPos, x, y)
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
