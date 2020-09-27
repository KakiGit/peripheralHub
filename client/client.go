package client

import (
	"fmt"
	"net"

	"github.com/kaki/peripheralHub/common"
)

func action(msg string) bool {
	return true
}

type Client struct {
	Address  common.Address
	Secret   common.Secret
	Platform common.Platform
}

func (client *Client) SyncWithServer(serverAddress common.Address, secretBytes []byte) {

	s, err := net.ResolveUDPAddr("udp", string(client.Address))
	c, err := net.DialUDP("udp", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()

	message := common.Message{
		SenderAddress:   common.AddressToBytes(client.Address),
		SenderPlatform:  client.Platform,
		ReceiverAddress: common.AddressToBytes(serverAddress),
		Event:           common.ServiceInit,
		EventEntity:     common.Client,
	}
	req := common.Encrypt(message, secretBytes)
	_, err = c.Write(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	input := common.Input{}
	input.Init()
	for {
		buffer := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp := common.Decrypt(buffer[0:n], secretBytes)
		fmt.Printf("Reply: %v\n", resp)
		if resp.SenderAddress == common.AddressToBytes(serverAddress) {
			input.InputFromClient(resp)
			client.Address = common.BytesToAddress(resp.ReceiverAddress)
		}
	}
}

func (client *Client) Init(serverAddress common.Address) {
	secretBytes := common.ReadSecret(client.Secret)

	client.SyncWithServer(serverAddress, secretBytes)
}
