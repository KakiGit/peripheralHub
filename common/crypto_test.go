package common

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCrypto(t *testing.T) {
	secret := "topSecret"
	encodedPwd := CreateEncodedKey(secret)
	fmt.Println(encodedPwd)
	key := ReadKey(encodedPwd)
	msgs := [][]byte{}
	for i := 32; i <= 128; i++ {
		msg := byte(i)
		msgs = append(msgs, []byte{msg})
	}
	for i := 0; i < 30; i++ {
		for _, msg := range msgs {
			encryptedMsg := Encrypt([]byte(msg), key)
			decryptedMsg := Decrypt(encryptedMsg, key)
			if res := bytes.Compare([]byte(msg), decryptedMsg); res != 0 {
				t.Log(res)
				t.Logf("secret: %s\nkey: %s\nmsg: %s\nencryptedMsg: %s\ndecrypted: %s\n",
					secret, key, string(msg), string(encryptedMsg), string(decryptedMsg))
				t.Failed()
			}
		}
	}
}
