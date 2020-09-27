// +build linux

package common

import "testing"

func TestOutput(t *testing.T) {
	com := make(chan InternalMsg, 100)
	go func() {
		for {
			<-com
		}
	}()
	OutputToServer(com)
}
