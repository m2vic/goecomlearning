package service

import (
	"context"
	"golearning/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) NewOrder(ctx context.Context, order domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}
func (m *MockOrderService) GetOrder(ctx context.Context, userId string) ([]domain.Order, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]domain.Order), nil
}
func (m *MockOrderService) UpdateOrderStatus(ctx context.Context, sessionId, status string) error {
	args := m.Called(ctx, sessionId)
	return args.Error(0)
}
