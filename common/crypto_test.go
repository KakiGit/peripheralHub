package common

import (
	"fmt"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCrypto(t *testing.T) {
	secret := "topSecret"
	encodedPwd := CreateEncodedSecret(secret)
	fmt.Println(encodedPwd)
	key := ReadSecret(encodedPwd)
	msgs := []Message{}
	for i := 0; i <= 255; i++ {
		for j := 0; j <= 255; j++ {
			nMsg := Message{
				SenderAddress:    [4]byte{192, 168, 0, 1},
				SenderPlatform:   Linux,
				ReceiverAddress:  [4]byte{192, 168, 0, 2},
				ReceiverPlatform: Windows,
				Event:            Event(j),
				EventEntity:      EventEntity(i),
			}
			msgs = append(msgs, nMsg)
		}
	}
	var wg sync.WaitGroup
	c := make(chan int, 10)
	for i := 0; i < 30; i++ {
		for _, msg := range msgs {
			c <- 1
			go func(wg *sync.WaitGroup, c chan int, msg Message, key []byte) {
				wg.Add(1)
				encryptedMsg := Encrypt(msg, key)
				decryptedMsg := Decrypt(encryptedMsg, key)
				if !cmp.Equal(decryptedMsg, msg) {
					t.Logf("secret: %s\nkey: %x\nmsg: %v\nencryptedMsg: %x\ndecrypted: %v\n",
						secret, key, msg, encryptedMsg, decryptedMsg)
					t.Failed()
				}
				<-c
				wg.Done()
			}(&wg, c, msg, key)
		}
		fmt.Println("running loop", i)
	}
	wg.Wait()
	fmt.Println("total of", len(msgs)*30)
}
