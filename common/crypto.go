package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/scrypt"
)

const nonceSize = 12
const saltSize = 32

func CreateEncodedKey(password string) string {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err.Error())
	}
	key, _ := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	encodedPwd := base64.StdEncoding.EncodeToString(key)
	return encodedPwd
}

func ReadKey(encodedPwd string) []byte {
	key, err := base64.StdEncoding.DecodeString(string(encodedPwd))
	if err != nil {
		panic(err.Error())
	}
	return key
}

func Encrypt(msg []byte, key []byte) []byte {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := aesgcm.Seal(nil, nonce, msg, nil)

	var encryptedMsg []byte
	encryptedMsg = append(encryptedMsg, nonce...)
	encryptedMsg = append(encryptedMsg, ciphertext...)
	ret := base64.StdEncoding.EncodeToString(encryptedMsg)
	return []byte(ret)
}

func Decrypt(encryptedMsg []byte, key []byte) []byte {
	msg, err := base64.StdEncoding.DecodeString(string(encryptedMsg))
	if err != nil {
		panic(err.Error())
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := msg[:nonceSize]
	decrptedMsg, err := aesgcm.Open(nil, nonce, msg[nonceSize:], nil)
	if err != nil {
		panic(err.Error())
	}
	return decrptedMsg
}
