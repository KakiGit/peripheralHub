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
