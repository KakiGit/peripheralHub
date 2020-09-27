package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"sort"
	"sync"

	"github.com/kaki/peripheralHub/common"
)

type Server struct {
	secret          []byte
	clients         map[string]string
	currentClient   net.Addr
	lock            *sync.Mutex
	serverAddresses []string
	platform        common.Platform
}

func (server *Server) serve(address string) {
	doneChan := make(chan error, 1)
	pc, err := net.ListenPacket("udp", address)
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
		maxBufferSize := 1024
		buffer := make([]byte, maxBufferSize)
		for {
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}
			msg := common.Decrypt(buffer[:n], server.secret)
			fmt.Printf("packet-received: bytes=%d from=%s text=%v\n",
				n, addr.String(), msg)
			receiverAddr := string(common.BytesToAddress(msg.ReceiverAddress))
			i := sort.SearchStrings(server.serverAddresses, receiverAddr)
			fmt.Println(i, len(server.serverAddresses), server.serverAddresses, receiverAddr)
			if i < len(server.serverAddresses) && server.serverAddresses[i] == receiverAddr {
				server.lock.Lock()
				if !server.hasClient(addr) {
					server.clients[addr.String()] = receiverAddr
				}
				// TODO: Change selection of currentClient to a better way
				server.currentClient = addr
				server.lock.Unlock()
			}
		}
	}(server)
	for {
		internalMsg := <-output.Com
		var receiverAddr common.Address
		var senderAddr [4]byte
		if server.currentClient != nil {
			receiverAddr = common.Address(server.currentClient.String())
			senderAddr = common.AddressToBytes(common.Address(server.findSenderAddr(server.currentClient)))
		} else {
			receiverAddr = ""
			senderAddr = [4]byte{}
		}
		msgToBeSent := common.Message{
			ReceiverAddress: common.AddressToBytes(receiverAddr),
			SenderAddress:   senderAddr,
			Event:           internalMsg.Event,
			EventEntity:     internalMsg.EventEntity,
			ExtraInfo:       internalMsg.ExtraInfo,
		}
		fmt.Printf("%v, \n", msgToBeSent)
		_, err := pc.WriteTo(common.Encrypt(msgToBeSent, server.secret), server.currentClient)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func ListenAndServe(address string, secret string) {
	sec := common.Secret(secret)
	encodedPwd := common.CreateEncodedSecret(sec)
	fmt.Println("Secret:", encodedPwd)
	d1 := []byte(encodedPwd)
	err := ioutil.WriteFile("/tmp/dat1", d1, 0644)
	fmt.Println(err)
	server := Server{
		secret:          common.ReadSecret(encodedPwd),
		lock:            &sync.Mutex{},
		serverAddresses: getServerAddresses(),
		clients:         make(map[string]string),
	}
	server.serve(address)
}

func getServerAddresses() []string {
	var svAddrs []string
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			svAddrs = append(svAddrs, ip.String())
			fmt.Println(ip)
		}
	}
	sort.Strings(svAddrs)
	return svAddrs
}

func (server *Server) hasClient(addr net.Addr) bool {
	for clientA := range server.clients {
		if clientA == addr.String() {
			return true
		}
	}
	return false
}

func (server *Server) findSenderAddr(addr net.Addr) string {
	senderAddr, ok := server.clients[addr.String()]
	if ok {
		return senderAddr
	} else {
		return ""
	}
}
