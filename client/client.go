package client

import (
	"fmt"
	"net"
	"time"

	"github.com/KakiGit/peripheralHub/common"
)

func action(msg string) bool {
	return true
}

type Client struct {
	Address   common.Address
	ClientId  string
	ServerId  string
	Port      string
	rsaCrypto common.RSACrypto
	Network   string
}

func (client *Client) SyncWithServer(serverAddress common.Address) {

	fullServerAddress := net.JoinHostPort(string(serverAddress), client.Port)
	conn, err := net.Dial(client.Network, fullServerAddress)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Printf("The TCP connection is %s\n", conn.LocalAddr().String())

	input := common.Input{}
	input.Init()
	buffer := make([]byte, 1024)
	// heartbeat
	go func(c net.Conn) {
		for {
			message := common.Message{
				SenderAddress:   common.AddressToBytes(client.Address),
				SenderId:        client.ClientId,
				ReceiverAddress: common.AddressToBytes(serverAddress),
				ReceiverId:      client.ServerId,
				Event:           common.ServiceHeartBeat,
				EventEntity:     common.Client,
			}
			msg := client.rsaCrypto.Encrypt(message)
			if err != nil {
				fmt.Println(err)
			}
			c.Write(msg)
			time.Sleep(60 * time.Second)
		}
	}(conn)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			// fmt.Println(err)
			continue
		}
		resp := client.rsaCrypto.Decrypt(buffer[0:n])
		// fmt.Printf("Resp: %v , %v, %v\n", resp, resp.SenderAddress, common.AddressToBytes(serverAddress))
		// if resp.SenderAddress == common.AddressToBytes(serverAddress) {
		input.InputFromClient(resp)
		// client.Address = common.BytesToAddress(resp.ReceiverAddress)
		// }
	}
}

func (client *Client) Init(serverAddress common.Address) {

	client.SyncWithServer(serverAddress)
}
