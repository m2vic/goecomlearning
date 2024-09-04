package service

import "github.com/stretchr/testify/mock"

type MockPasswordGenerator struct {
	mock.Mock
}

func (m *MockPasswordGenerator) RandomPassword() (string, error) {
	args := m.Called()
	return args.String(0), nil
}
