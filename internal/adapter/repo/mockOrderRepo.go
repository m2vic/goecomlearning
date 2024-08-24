package repo

import (
	"context"
	"golearning/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) NewOrder(ctx context.Context, order domain.Order) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockOrderRepo) GetOrder(ctx context.Context, userId string) ([]domain.Order, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]domain.Order), nil
}
func (m *MockOrderRepo) UpdateOrderStatus(ctx context.Context, sessionId, status string) error {
	args := m.Called(ctx, sessionId, status)
	return args.Error(0)
}
