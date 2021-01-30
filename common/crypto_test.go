package common

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/ssh"
)

func TestAES(t *testing.T) {
	secret := Secret("topSecret")
	aesCrypto := AESCrypto{}
	encodedPwd := aesCrypto.CreateEncodedSecret(secret)
	fmt.Println(encodedPwd)
	aesCrypto.ReadSecret(encodedPwd)
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
				encryptedMsg := aesCrypto.Encrypt(msg)
				decryptedMsg := aesCrypto.Decrypt(encryptedMsg)
				if !cmp.Equal(decryptedMsg, msg) {
					t.Logf("secret: %s\nkey: %x\nmsg: %v\nencryptedMsg: %x\ndecrypted: %v\n",
						secret, key, msg, encryptedMsg, decryptedMsg)
					t.Failed()
				}
				<-c
				wg.Done()
			}(&wg, c, msg, aesCrypto.Key)
		}
		fmt.Println("running loop", i)
	}
	wg.Wait()
	fmt.Println("total of", len(msgs)*30)
}

func TestRSA(t *testing.T) {

	clientRsaCrypto := RSACrypto{}
	clientRsaCrypto.ReadSecret("/home/kaki/.ssh/id_rsa_test_client.pub", "/home/kaki/.ssh/id_rsa_test_client")
	publicKeyData, err := ioutil.ReadFile("/home/kaki/.ssh/id_rsa_test_server.pub")
	if err != nil {
		fmt.Println(err)
	}
	parsed, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyData)
	if err != nil {
		fmt.Println(err)
	}
	publicKey := parsed.(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey)
	clientRsaCrypto.targetPublicKey = (RSAPubKey)(*publicKey)

	serverRsaCrypto := RSACrypto{}
	serverRsaCrypto.ReadSecret("/home/kaki/.ssh/id_rsa_test_server.pub", "/home/kaki/.ssh/id_rsa_test_server")
	publicKeyData, err = ioutil.ReadFile("/home/kaki/.ssh/id_rsa_test_client.pub")
	if err != nil {
		fmt.Println(err)
	}
	parsed, _, _, _, err = ssh.ParseAuthorizedKey(publicKeyData)
	if err != nil {
		fmt.Println(err)
	}
	publicKey = parsed.(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey)
	serverRsaCrypto.targetPublicKey = (RSAPubKey)(*publicKey)

	msgs := []Message{}
	for i := 0; i <= 9; i++ {
		for j := 0; j <= 123; j++ {
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
		counter := 0
		for _, msg := range msgs {
			c <- counter
			go func(wg *sync.WaitGroup, c chan int, msg Message) {
				wg.Add(1)
				encryptedMsg := serverRsaCrypto.Encrypt(msg)
				decryptedMsg := clientRsaCrypto.Decrypt(encryptedMsg)
				if !cmp.Equal(decryptedMsg, msg) {
					t.Logf("msg: %v\nencryptedMsg: %x\ndecrypted: %v\n",
						msg, encryptedMsg, decryptedMsg)
					t.Failed()
				}
				<-c
				// fmt.Println(current)
				counter += 1
				wg.Done()
			}(&wg, c, msg)
		}
		fmt.Println("running loop", i)
	}
	wg.Wait()
	fmt.Println("total of", len(msgs)*30)
}
