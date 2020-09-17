package server

import (
	"fmt"
	"net"
	"time"

	"github.com/kaki/peripheralHub/common"
)

func server(address string) {
	to := 1 * time.Second
	timeout := &to
	maxBufferSize := 1024
	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pc.Close()
	fmt.Println(pc.LocalAddr().String())
	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	for {
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("packet-received: bytes=%d from=%s text=%s\n",
			n, addr.String(), common.Decrypt(buffer[:n]))
		deadline := time.Now().Add(*timeout)
		err = pc.SetWriteDeadline(deadline)
		if err != nil {
			doneChan <- err
			return
		}

		keyToBeSent := string(buffer[:n])
		n, err = pc.WriteTo(common.Encrypt(keyToBeSent), addr)
		if err != nil {
			doneChan <- err
			return
		}
	}
}

func ListenAndServe(address string) {
	server(address)
}
