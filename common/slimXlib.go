// +build linux

package common

/*
 #cgo LDFLAGS: -L/usr/include/X11 -lX11 -lXtst -lXext
 #include <X11/Xlib.h>
 #include <X11/Intrinsic.h>
 #include <X11/extensions/XTest.h>
 #include <X11/XKBlib.h>
int mDefaultScreen(Display *display)
{
	return DefaultScreen(display);
}
Window mRootWindow(Display* display, int s)
{
	return RootWindow(display, s);
}
unsigned int getType(XEvent *event) {
	return event->type;
}
unsigned int getKeyCode(XEvent *event) {
	return event->xkey.keycode;
}
void getCursorPosition(XEvent *event, int *x, int *y) {
	*x = event->xkey.x_root;
	*y = event->xkey.y_root;
}
*/
import "C"

type Display C.Display
type Screen C.Screen
type Window C.Window
type XEvent C.XEvent

const (
	KeyPress          = C.KeyPress
	KeyRelease        = C.KeyRelease
	ButtonPress       = C.ButtonPress
	ButtonRelease     = C.ButtonRelease
	MotionNotify      = C.MotionNotify
	ButtonPressMask   = C.ButtonPressMask
	ButtonReleaseMask = C.ButtonReleaseMask
	PointerMotionMask = C.PointerMotionMask
	GrabModeAsync     = C.GrabModeAsync
	None              = C.None
	CurrentTime       = C.CurrentTime
	True              = C.True
	False             = C.False
)

func XTestGrabControl(display *Display, onOff int) {
	if onOff == False {
		C.XSync((*C.Display)(display), False)
	}
	C.XTestGrabControl((*C.Display)(display), C.int(onOff))
}

func FakeKeyEvent(display *Display, keycode uint, eventType int, delay int) {
	var eventtype C.int
	switch eventType {
	case KeyPress:
		eventtype = True
	case KeyRelease:
		eventtype = False
	}
	C.XTestFakeKeyEvent((*C.Display)(display), C.uint(keycode), eventtype, 0)
}

func XKeycodeToKeysym(display *Display, keycode uint) uint64 {
	return uint64(C.XkbKeycodeToKeysym((*C.Display)(display), C.uchar(keycode), 0, 0))
}

func XKeysymToString(keysym uint64) string {
	return C.GoString((*C.char)(C.XKeysymToString(C.KeySym(keysym))))
}

func XStringToKeysym(keyString string) uint64 {
	return uint64(C.XStringToKeysym(C.CString(keyString)))
}

func XKeysymToKeycode(display *Display, keysym uint64) uint {
	return uint(C.XKeysymToKeycode((*C.Display)(display), C.KeySym(keysym)))
}

func XOpenDisplay() *Display {
	return (*Display)(C.XOpenDisplay(nil))
}

func XDefaultRootWindow(display *Display) Window {
	return (Window)(C.XDefaultRootWindow((*C.Display)(display)))
}

func XGrabPointer(display *Display, window Window) {
	C.XGrabPointer((*C.Display)(display), (C.Window)(window), 0, ButtonPressMask|ButtonReleaseMask|
		PointerMotionMask, GrabModeAsync, GrabModeAsync, None,
		None, CurrentTime)
}

func XGrabKeyboard(display *Display, window Window) {
	C.XGrabKeyboard((*C.Display)(display), (C.Window)(window), 0, GrabModeAsync, GrabModeAsync, CurrentTime)
}

func XNextEvent(display *Display, event *XEvent) {
	C.XNextEvent((*C.Display)(display), (*C.XEvent)(event))
}

func GetXEventType(event *XEvent) int {
	return int(C.getType((*C.XEvent)(event)))
}

func GetKeyCode(event *XEvent) uint {
	return uint(C.getKeyCode((*C.XEvent)(event)))
}

func GetCursorPosition(event *XEvent) (int, int) {
	var x, y C.int
	C.getCursorPosition((*C.XEvent)(event), &x, &y)
	return int(x), int(y)
}
