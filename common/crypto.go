package common

func Encrypt(msg string) []byte {
	return []byte(msg)
}

func Decrypt(encryptedKey []byte) string {
	return string(encryptedKey)
}
