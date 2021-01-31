package server

import (
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/KakiGit/peripheralHub/common"
)

func TestServer(t *testing.T) {

	server := Server{
		Address:         "",
		Port:            "",
		lock:            &sync.Mutex{},
		serverAddresses: getServerAddresses(),
		clients:         make(map[string]net.Conn),
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	sshDir := filepath.Join(homeDir, ".ssh")
	publicKeyPath := filepath.Join(sshDir, "id_rsa_test.pub")
	privateKeyPath := filepath.Join(sshDir, "id_rsa_test")
	authorizedKeysPath := filepath.Join(sshDir, "authorized_keys")
	server.rsaCrypto.ReadSecret(publicKeyPath, privateKeyPath)
	server.clientPublicKeys, _ = common.ReadAuthorizedKeys(authorizedKeysPath)
	server.Serve()
}
