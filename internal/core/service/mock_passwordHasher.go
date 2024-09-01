package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) ComparePassword(ctx context.Context, password, hash string) error {
	args := m.Called(ctx, password, hash)
	return args.Error(0)
}

func (m *MockPasswordHasher) HashPassword(ctx context.Context, password string) ([]byte, error) {
	args := m.Called(ctx, password)
	return args.Get(0).([]byte), nil
}
