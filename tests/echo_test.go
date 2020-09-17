package tests

import (
	"testing"
	"time"

	"github.com/kaki/peripheralHub/client"
	"github.com/kaki/peripheralHub/common"
	"github.com/kaki/peripheralHub/server"
)

func TestEcho(t *testing.T) {
	sender := common.Sender{Address: "127.0.0.1:9878"}
	receiver := common.Receiver{Secret: "receiverSecret"}
	go server.ListenAndServe(sender.Address)
	time.Sleep(100 * time.Millisecond)
	client.SyncWithServer(sender, receiver)
}
