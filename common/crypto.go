package common

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/gob"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh"
)

const nonceSize = 12
const saltSize = 32

type AESCrypto struct {
	Key AESKey
}

type RSACrypto struct {
	publicKey        RSAPubKey
	PublicKeyComment string
	privateKey       RSAPriKey
	TargetPublicKey  RSAPubKey
}

func (aesCrypto *AESCrypto) CreateEncodedSecret(password Secret) Secret {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		fmt.Println(err)
	}
	key, _ := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	encodedPwd := base64.StdEncoding.EncodeToString(key)
	return Secret(encodedPwd)
}

func (rsaCrypto *RSACrypto) CreateEncodedSecret() {

}

func (aesCrypto *AESCrypto) ReadSecret(encodedPwd Secret) {
	key, err := base64.StdEncoding.DecodeString(string(encodedPwd))
	if err != nil {
		fmt.Println(err)
	}
	aesCrypto.Key = key
}

func (rsaCrypto *RSACrypto) ReadSecret(publicKeyPath string, privateKeyPath string) {
	publicKeyData, err := ioutil.ReadFile(publicKeyPath)

	parsed, comment, _, _, err := ssh.ParseAuthorizedKey(publicKeyData)
	CheckError(err)
	publicKey := parsed.(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey)
	rsaCrypto.publicKey = (RSAPubKey)(*publicKey)
	rsaCrypto.PublicKeyComment = comment

	privateKeyData, err := ioutil.ReadFile(privateKeyPath)
	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		panic(errors.New("ssh: no key found"))
	}
	privateKeyParsed, err := ssh.ParseRawPrivateKey(privateKeyData)
	switch block.Type {
	case "OPENSSH PRIVATE KEY":
	case "RSA PRIVATE KEY":
		privateKeyParsed, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		CheckError(err)
	default:
		panic(errors.New("Unsupported key type"))
	}
	privateKey := privateKeyParsed.(*rsa.PrivateKey)
	rsaCrypto.privateKey = (RSAPriKey)(*privateKey)
}

func ReadAuthorizedKeys(autorizedKeysPath string) (map[string]RSAPubKey, []string) {
	comments := []string{}
	authorizedKeys := make(map[string]RSAPubKey)
	data, err := ioutil.ReadFile(autorizedKeysPath)
	CheckError(err)
	var parsed ssh.PublicKey
	var comment string
	for {
		parsed, comment, _, data, err = ssh.ParseAuthorizedKey(data)
		comments = append(comments, comment)
		CheckError(err)
		publicKey := parsed.(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey)
		authorizedKeys[comment] = (RSAPubKey)(*publicKey)
		if len(data) == 0 {
			break
		}
	}

	return authorizedKeys, comments
}

func (aesCrypto *AESCrypto) Encrypt(msg Message) []byte {
	block, err := aes.NewCipher(aesCrypto.Key)
	if err != nil {
		fmt.Println(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}
	ciphertext := aesgcm.Seal(nil, nonce, EncodeMessage(msg), nil)

	var encryptedMsg []byte
	encryptedMsg = append(encryptedMsg, nonce...)
	encryptedMsg = append(encryptedMsg, ciphertext...)
	ret := base64.StdEncoding.EncodeToString(encryptedMsg)
	return []byte(ret)
}

func (rsaCrypto *RSACrypto) Encrypt(msg Message) []byte {
	rng := rand.Reader
	label := []byte("peripheralHub")

	encodedMsg, signiture := rsaCrypto.EncodeMessage(msg)
	encryptedMsg, err := rsa.EncryptOAEP(sha256.New(), rng, (*rsa.PublicKey)(&rsaCrypto.TargetPublicKey), encodedMsg, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		panic(err)
	}

	ret := base64.StdEncoding.EncodeToString(encryptedMsg) + "\n" + base64.StdEncoding.EncodeToString(signiture)
	return []byte(ret)
}

func (rsaCrypto *RSACrypto) EncodeMessage(msg Message) ([]byte, []byte) {
	rng := rand.Reader
	encodedMsg := EncodeMessage(msg)
	pssOptions := &rsa.PSSOptions{SaltLength: saltSize, Hash: crypto.SHA256}
	msgHash := pssOptions.Hash.HashFunc().New()
	_, err := msgHash.Write(encodedMsg)
	CheckError(err)
	msgHashSum := msgHash.Sum(nil)
	signiture, err := rsa.SignPSS(rng, (*rsa.PrivateKey)(&rsaCrypto.privateKey), crypto.SHA256, msgHashSum, pssOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
	}

	return encodedMsg, signiture
}

func (aesCrypto *AESCrypto) Decrypt(encryptedMsg []byte) Message {
	msg, err := base64.StdEncoding.DecodeString(string(encryptedMsg))
	if err != nil {
		fmt.Println(err)
	}
	block, err := aes.NewCipher(aesCrypto.Key)
	if err != nil {
		fmt.Println(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
	}
	nonce := msg[:nonceSize]
	decrptedMsg, err := aesgcm.Open(nil, nonce, msg[nonceSize:], nil)
	if err != nil {
		fmt.Println(err)
	}

	return DecodeMessage(decrptedMsg)
}

func (rsaCrypto *RSACrypto) Decrypt(msg []byte) Message {
	msgs := strings.Split(string(msg), "\n")
	encryptedMsg, err := base64.StdEncoding.DecodeString(msgs[0])
	if err != nil {
		fmt.Println(err)
	}
	signiture, err := base64.StdEncoding.DecodeString(msgs[1])
	if err != nil {
		fmt.Println(err)
	}

	label := []byte("peripheralHub")
	rng := rand.Reader
	decryptedMsg, err := rsa.DecryptOAEP(sha256.New(), rng, (*rsa.PrivateKey)(&rsaCrypto.privateKey), encryptedMsg, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from decryption: %s\n", err)
		panic(err)
	}
	return rsaCrypto.DecodeMessage(decryptedMsg, signiture)
}

func (rsaCrypto *RSACrypto) DecodeMessage(encodedMsg []byte, signiture []byte) Message {

	pssOptions := &rsa.PSSOptions{SaltLength: saltSize, Hash: crypto.SHA256}
	msgHash := pssOptions.Hash.HashFunc().New()
	_, err := msgHash.Write(encodedMsg)
	if err != nil {
		fmt.Println(err)
	}
	msgHashSum := msgHash.Sum(nil)
	err = rsa.VerifyPSS((*rsa.PublicKey)(&rsaCrypto.TargetPublicKey), crypto.SHA256, msgHashSum, signiture, pssOptions)
	if err != nil {
		fmt.Println(err)
		return Message{}
	}
	return DecodeMessage(encodedMsg)
}

func EncodeMessage(msg Message) []byte {
	var encodedMsg bytes.Buffer
	enc := gob.NewEncoder(&encodedMsg)
	err := enc.Encode(msg)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return encodedMsg.Bytes()
}

func DecodeMessage(msg []byte) Message {
	dec := gob.NewDecoder(bytes.NewReader(msg))
	var decodedMsg Message
	err := dec.Decode(&decodedMsg)
	if err != nil {
		fmt.Println(err)
		return Message{}
	}
	return decodedMsg
}
