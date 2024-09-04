package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type CryptoService struct {
	EncryptionKey []byte
}

func NewCryptoService(key []byte) *CryptoService {
	return &CryptoService{EncryptionKey: key}
}
func (s *CryptoService) Encrypt(email string) (string, error) {
	block, err := aes.NewCipher(s.EncryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(email))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(email))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}
func (s *CryptoService) Decrypt(email string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(email)

	block, err := aes.NewCipher(s.EncryptionKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
