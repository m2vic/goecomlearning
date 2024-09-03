package service

import (
	"crypto/rand"
	"log"
)

type PasswordGenerator struct{}

func (g *PasswordGenerator) RandomPassword() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	charsetLen := len(charset)
	password := make([]byte, 12)
	if _, err := rand.Read(password); err != nil {
		log.Fatal(err)
	}
	for i := range password {
		password[i] = charset[int(password[i])%charsetLen]
	}
	return string(password), nil
}
