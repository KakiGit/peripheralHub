package common

type Sender struct {
	Address         string
	Devices         []Device
	Receivers       []*Receiver
	CurrentReceiver *Receiver
	Secret          Secret
}
