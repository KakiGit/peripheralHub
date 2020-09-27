package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

const nonceSize = 12
const saltSize = 32

func CreateEncodedSecret(password Secret) Secret {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err.Error())
	}
	key, _ := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	encodedPwd := base64.StdEncoding.EncodeToString(key)
	return Secret(encodedPwd)
}

func ReadSecret(encodedPwd Secret) []byte {
	key, err := base64.StdEncoding.DecodeString(string(encodedPwd))
	if err != nil {
		panic(err.Error())
	}
	return key
}

func Encrypt(msg Message, key []byte) []byte {

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
	ciphertext := aesgcm.Seal(nil, nonce, EncodeMessage(msg), nil)

	var encryptedMsg []byte
	encryptedMsg = append(encryptedMsg, nonce...)
	encryptedMsg = append(encryptedMsg, ciphertext...)
	ret := base64.StdEncoding.EncodeToString(encryptedMsg)
	return []byte(ret)
}

func Decrypt(encryptedMsg []byte, key []byte) Message {
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

	return DecodeMessage(decrptedMsg)
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
