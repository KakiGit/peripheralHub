// +build linux

package common

import "testing"

func TestOutput(t *testing.T) {
	output := Output{}
	com := make(chan InternalMsg, 100)
	go func() {
		for {
			<-com
		}
	}()
	output.OutputToServer(com)
}
