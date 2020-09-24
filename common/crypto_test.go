package common

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCrypto(t *testing.T) {
	secret := "topSecret"
	encodedPwd := CreateEncodedKey(secret)
	fmt.Println(encodedPwd)
	key := ReadKey(encodedPwd)
	msgs := []Message{}
	for i := 0; i <= 255; i++ {
		for j := 0; j <= 255; j++ {
			nMsg := Message{
				Sender:      [4]byte{192, 168, 0, 1},
				Receiver:    [4]byte{192, 168, 0, 2},
				Event:       Event(j),
				EventEntity: EventEntity(i),
			}
			msgs = append(msgs, nMsg)
		}
	}
	for i := 0; i < 30; i++ {
		for _, msg := range msgs {
			encryptedMsg := Encrypt(msg, key)
			decryptedMsg := Decrypt(encryptedMsg, key)
			if !cmp.Equal(decryptedMsg, msg) {
				t.Logf("secret: %s\nkey: %x\nmsg: %v\nencryptedMsg: %x\ndecrypted: %v\n",
					secret, key, msg, encryptedMsg, decryptedMsg)
				t.Failed()
			}
		}
	}
}
