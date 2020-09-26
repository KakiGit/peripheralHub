// +build linux

package common

import (
	"testing"
)

func TestSender(t *testing.T) {
	com := make(chan InternalMsg, 100)
	go func() {
		for {
			<-com
		}
	}()
	ListenXEvent(com)
}
