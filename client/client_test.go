package client

import (
	"testing"

	"github.com/kaki/peripheralHub/common"
)

func TestClient(t *testing.T) {
	client := Client{
		Secret: "",
	}
	client.Init(common.Address(""))
}
