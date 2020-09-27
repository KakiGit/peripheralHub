package server

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/kaki/peripheralHub/common"
)

type Server struct {
	secret        []byte
	clients       []net.Addr
	currentClient net.Addr
	lock          *sync.Mutex
}

func serve(address string, secret string) {
	server := Server{
		secret: common.ReadSecret(common.CreateEncodedSecret(common.Secret(secret))),
		lock:   &sync.Mutex{},
	}
	doneChan := make(chan error, 1)
	to := 1 * time.Second
	timeout := &to
	maxBufferSize := 1024
	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	deadline := time.Now().Add(*timeout)
	err = pc.SetWriteDeadline(deadline)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pc.Close()
	fmt.Println(pc.LocalAddr().String())
	output := common.Output{}
	output.Init()
	go output.OutputToServer()
	go func(server *Server) {
		buffer := make([]byte, maxBufferSize)
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			doneChan <- err
			return
		}
		msg := common.Decrypt(buffer[:n], server.secret)
		fmt.Printf("packet-received: bytes=%d from=%s text=%v\n",
			n, addr.String(), msg)
		addrPart := strings.Split(addr.String(), ":")[0]
		if common.Address(addrPart) == common.BytesToAddress(msg.SenderAddress) {
			server.lock.Lock()
			server.clients = append(server.clients, addr)
			server.currentClient = addr
			server.lock.Unlock()
		}
	}(&server)
	for {

		internalMsg := <-output.Com
		msgToBeSent := common.Message{
			Event:       internalMsg.Event,
			EventEntity: internalMsg.EventEntity,
			ExtraInfo:   internalMsg.ExtraInfo,
		}
		fmt.Printf("%v\n", msgToBeSent)
		pc.WriteTo(common.Encrypt(msgToBeSent, server.secret), server.currentClient)
	}
}

func ListenAndServe(address string, secret string) {
	serve(address, secret)
}
