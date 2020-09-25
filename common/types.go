package common

type Address string
type Secret string
type Device string
type Event byte
type EventEntity byte
type Message struct {
	Sender      [4]byte
	Receiver    [4]byte
	Event       Event
	EventEntity EventEntity
}

const (
	KeyboardButtonDown Event = iota
	KeyboardButtonUp
	MouseButtonDown
	MouseButtonUp
	MouseButtonClick
	MouseRelativeMove
	MouseAbsoluteMove
	MouseWheelScroll
)

// ISO keyboard layout is implemented for now
const (
	MouseLeftButton EventEntity = iota
	MouseRightButton
	MouseMiddleButton
	MouseOptButton1
	MouseOptButton2
	MouseWheel
	MouseCursor
	KeyVolumeMute
	KeyVolumeUp
	KeyVolumeDown
	KeyMediaStop
	KeyMediaPause
	KeyMediaNext
	KeyMediaPrevious
	KeyEsc
	KeyWinCmd
	KeyLWinCmd
	KeyRWinCmd
	KeyAlt
	KeyLAlt
	KeyRAlt
	KeyCtrl
	KeyLCtrl
	KeyRCtrl
	KeyShift
	KeyLShift
	KeyRShift
	KeyTab
	KeyCaps
	KeyAltGr
	KeyEnter
	KeySpace
	KeyBackspace
	KeyDelete
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyInsert
	KeyPrintScreen
	KeyScrollLock
	KeyPause
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyTilde
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyLeftBracket
	KeyBar
	KeyRightBracket
	KeyNumLock
	KeyNumPad0
	KeyNumPad1
	KeyNumPad2
	KeyNumPad3
	KeyNumPad4
	KeyNumPad5
	KeyNumPad6
	KeyNumPad7
	KeyNumPad8
	KeyNumPad9
	KeyNumPadEnter
	KeyNumPadMul
	KeyNumPadDiv
	KeyNumPadAdd
	KeyNumPadSub
	KeyNumPadDecimal
)
