package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct{}

func (s *PasswordHasher) ComparePassword(ctx context.Context, password, hash string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password invalid")
	}
	return nil
}

func (s *PasswordHasher) HashPassword(ctx context.Context, password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("fail to generate password")
	}
	return hash, nil
}
