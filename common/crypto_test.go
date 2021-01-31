package common

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	sshDir := filepath.Join(homeDir, ".ssh")
	clientPublicKeyPath := filepath.Join(sshDir, "id_rsa_test_client.pub")
	clientPrivateKeyPath := filepath.Join(sshDir, "id_rsa_test_client")
	serverPublicKeyPath := filepath.Join(sshDir, "id_rsa_test_server.pub")
	serverPrivateKeyPath := filepath.Join(sshDir, "id_rsa_test_server")

	clientRsaCrypto := RSACrypto{}
	clientRsaCrypto.ReadSecret(clientPublicKeyPath, clientPrivateKeyPath)
	publicKey, serverName := ReadAuthorizedKeys(serverPublicKeyPath)
	clientRsaCrypto.TargetPublicKey = publicKey[serverName[0]]

	serverRsaCrypto := RSACrypto{}
	serverRsaCrypto.ReadSecret(serverPublicKeyPath, serverPrivateKeyPath)
	publicKey, clientName := ReadAuthorizedKeys(clientPublicKeyPath)
	serverRsaCrypto.TargetPublicKey = publicKey[clientName[0]]

	msgs := []Message{}
	for i := 0; i <= 9; i++ {
		for j := 0; j <= 123; j++ {
			nMsg := Message{
				SenderAddress:    [4]byte{192, 168, 0, 1},
				SenderId:         serverName[0],
				SenderPlatform:   Linux,
				ReceiverAddress:  [4]byte{192, 168, 0, 2},
				ReceiverId:       clientName[0],
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
				fmt.Println(string(encryptedMsg))
				decryptedMsg := clientRsaCrypto.Decrypt(encryptedMsg)
				fmt.Println(decryptedMsg)
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
