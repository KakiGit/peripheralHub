package client

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KakiGit/peripheralHub/common"
)

func TestClient(t *testing.T) {

	client := Client{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	sshDir := filepath.Join(homeDir, ".ssh")
	publicKeyPath := filepath.Join(sshDir, "id_rsa_test.pub")
	privateKeyPath := filepath.Join(sshDir, "id_rsa_test")
	authorizedKeysPath := filepath.Join(sshDir, "authorized_keys")
	client.rsaCrypto.ReadSecret(publicKeyPath, privateKeyPath)
	authorizedKeys, comments := common.ReadAuthorizedKeys(authorizedKeysPath)
	client.rsaCrypto.TargetPublicKey = authorizedKeys[comments[0]]
	client.ClientId = client.rsaCrypto.PublicKeyComment
	client.ServerId = comments[0]

	client.Init(common.Address(""))
}
