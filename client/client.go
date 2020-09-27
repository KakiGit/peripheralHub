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

	s, err := net.ResolveUDPAddr("udp", string(serverAddress)+":9900")
	pc, err := net.ListenPacket("udp", string(client.Address)+":9900")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pc.Close()

	fmt.Printf("The UDP server is %s\n", string(client.Address)+":9900")

	message := common.Message{
		SenderAddress:   common.AddressToBytes(client.Address),
		SenderPlatform:  client.Platform,
		ReceiverAddress: common.AddressToBytes(serverAddress),
		Event:           common.ServiceInit,
		EventEntity:     common.Client,
	}
	req := common.Encrypt(message, secretBytes)
	_, err = pc.WriteTo(req, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	input := common.Input{}
	input.Init()
	buffer := make([]byte, 1024)
	for {
		n, _, err := pc.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp := common.Decrypt(buffer[0:n], secretBytes)
		fmt.Printf("Resp: %v , %v, %v\n", resp, resp.SenderAddress, common.AddressToBytes(serverAddress))
		// if resp.SenderAddress == common.AddressToBytes(serverAddress) {
		input.InputFromClient(resp)
		// client.Address = common.BytesToAddress(resp.ReceiverAddress)
		// }
	}
}

func (client *Client) Init(serverAddress common.Address) {
	secretBytes := common.ReadSecret(client.Secret)

	client.SyncWithServer(serverAddress, secretBytes)
}
