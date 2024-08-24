package repo

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockProductCached struct {
	mock.Mock
}

func (m *MockProductCached) SetProduct(ctx context.Context, json []byte) error {
	return nil
}

func (m *MockProductCached) GetProduct(context.Context) (string, error) {
	return "", nil
}
