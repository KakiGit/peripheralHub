package client

import (
	"fmt"
	"net"
	"time"

	"github.com/kaki/peripheralHub/common"
)

func action(msg string) bool {
	return true
}

func SyncWithServer(sender common.Sender, receiver common.Receiver) {

	s, err := net.ResolveUDPAddr("udp", sender.Address)
	c, err := net.DialUDP("udp", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()

	timeSpan := 500 * time.Millisecond
	end := time.Now().Add(timeSpan)

	for time.Now().Before(end) {

		text := receiver.Secret
		data := common.Encrypt(string(text))
		_, err = c.Write(data)

		if err != nil {
			fmt.Println(err)
			return
		}

		buffer := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Reply: %s\n", common.Decrypt(buffer[0:n]))
		time.Sleep(10 * time.Millisecond)
	}
}

func Init() {
	sender := common.Sender{Address: "127.0.0.1:9878"}
	receiver := common.Receiver{Secret: "receiverSecret"}

	SyncWithServer(sender, receiver)
}
