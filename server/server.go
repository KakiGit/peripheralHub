package server

import (
	"fmt"
	"net"
	"sort"
	"sync"

	"github.com/KakiGit/peripheralHub/common"
)

type Server struct {
	Address          common.Address
	ServerId         string
	Port             string
	clients          map[string]net.Conn
	currentClient    net.Conn
	currentClientId  string
	clientPublicKeys map[string]common.RSAPubKey
	rsaCrypto        common.RSACrypto
	lock             *sync.Mutex
	serverAddresses  []string
	Network          string
}

func (server *Server) Serve() {
	doneChan := make(chan error, 1)
	fullAddr := net.JoinHostPort(string(server.Address), server.Port)
	ln, err := net.Listen(server.Network, fullAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()
	fmt.Println(ln.Addr().String())
	output := common.Output{}
	output.Init()
	go output.OutputToServer()
	go func(server *Server) {
		maxBufferSize := 1024
		buffer := make([]byte, maxBufferSize)
		for {
			conn, err := ln.Accept()
			if err != nil {
				doneChan <- err
			}
			go func(c net.Conn) {
				for {
					n, err := conn.Read(buffer)
					if err != nil {
						fmt.Println(err)
					}
					msg := server.rsaCrypto.Decrypt(buffer[:n])
					fmt.Printf("data received: bytes=%d from=%s text=%v\n",
						n, conn.RemoteAddr().String(), msg)
					receiverAddr := common.BytesToAddress(msg.ReceiverAddress)
					if server.Address == receiverAddr {
						server.lock.Lock()
						if !server.hasClient(msg.SenderId) {
							server.clients[msg.SenderId] = conn
							defer delete(server.clients, msg.SenderId)
						}
						server.currentClient = conn
						server.currentClientId = msg.SenderId
						server.rsaCrypto.TargetPublicKey = server.clientPublicKeys[msg.SenderId]
						server.lock.Unlock()
					}
				}
			}(conn)
		}
	}(server)
	for {
		internalMsg := <-output.Com
		var receiverAddr common.Address
		var currentClientId string
		server.lock.Lock()
		if server.currentClient != nil {
			receiverAddr = common.Address(server.currentClient.RemoteAddr().String())
			currentClientId = server.currentClientId
			server.lock.Unlock()
		} else {
			server.lock.Unlock()
			continue
		}
		msgToBeSent := common.Message{
			ReceiverAddress: common.AddressToBytes(receiverAddr),
			ReceiverId:      currentClientId,
			SenderAddress:   common.AddressToBytes(server.Address),
			SenderId:        server.ServerId,
			Event:           internalMsg.Event,
			EventEntity:     internalMsg.EventEntity,
			ExtraInfo:       internalMsg.ExtraInfo,
		}
		fmt.Printf("%v, \n", msgToBeSent)
		_, err := server.currentClient.Write(server.rsaCrypto.Encrypt(msgToBeSent))
		if err != nil {
			fmt.Println(err)
		}
	}
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
				fmt.Println("ipnet", ip)
			case *net.IPAddr:
				ip = v.IP
				fmt.Println("ipaddr", ip)
			}
			svAddrs = append(svAddrs, ip.String())
		}
	}
	sort.Strings(svAddrs)
	return svAddrs
}

func (server *Server) hasClient(clientId string) bool {
	_, ok := server.clients[clientId]
	if ok {
		return true
	}

	return false
}
