// +build linux

package common

import "testing"

func TestOutput(t *testing.T) {
	output := Output{}
	go func(com chan InternalMsg) {
		for {
			<-com
		}
	}(output.Com)
	output.OutputToServer()
}
