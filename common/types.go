package common

import (
	"strconv"
	"strings"
)

type Address string
type Secret string
type Device string
type Platform int
type Event byte
type EventEntity byte

type Message struct {
	SenderAddress    [4]byte
	SenderPlatform   Platform
	ReceiverAddress  [4]byte
	ReceiverPlatform Platform
	Event            Event
	EventEntity      EventEntity
	ExtraInfo        [4]int
}

func AddressToBytes(address Address) [4]byte {
	tmp := string(address)
	if strings.Contains(tmp, ":") {
		tmp = strings.Split(tmp, ":")[0]
	}
	parts := strings.Split(tmp, ".")
	ret := [4]byte{}
	for i, part := range parts {
		pt, _ := strconv.Atoi(part)
		ret[i] = byte(pt)
	}
	return ret
}

func BytesToAddress(addressBytes [4]byte) Address {
	var addr string
	for _, partByte := range addressBytes {
		part := strconv.Itoa(int(partByte))
		addr += part + "."
	}
	return Address(strings.Trim(addr, "."))
}

type InternalMsg struct {
	EventEntity EventEntity
	Event       Event
	ExtraInfo   [4]int
}

const (
	Linux Platform = iota
	Windows
	Darwin
)

const (
	ButtonDown Event = iota
	ButtonUp
	ButtonClick
	MouseRelativeMove
	MouseAbsoluteMove
	MouseWheelScrollUp
	MouseWheelScrollDown
	ServiceInit
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
	Client
	Server
)
