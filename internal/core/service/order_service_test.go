package service_test

import (
	"context"
	"golearning/internal/adapter/repo"
	"golearning/internal/core/domain"
	"golearning/internal/core/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrder(t *testing.T) {
	mockOrderRepo := new(repo.MockOrderRepo)
	orderService := service.NewOrderService(mockOrderRepo)

	ctx := context.Background()
	userId := "1"
	var orders []domain.Order
	order := domain.Order{TotalPrice: 123, OrderID: "122"}
	orders = append(orders, order)
	expected := orders
	mockOrderRepo.On("GetOrder", ctx, userId).Return(expected)
	actual, err := orderService.GetOrder(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestNewOrder(t *testing.T) {
	mockOrderRepo := new(repo.MockOrderRepo)
	orderService := service.NewOrderService(mockOrderRepo)

	ctx := context.Background()

	order := domain.Order{TotalPrice: 123, OrderID: "122"}
	mockOrderRepo.On("NewOrder", ctx).Return(nil)
	actual := orderService.NewOrder(ctx, order)
	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}

func TestUpdateOrderStatus(t *testing.T) {
	mockOrderRepo := new(repo.MockOrderRepo)
	orderService := service.NewOrderService(mockOrderRepo)
	ctx := context.Background()
	sessionId := "1"
	status := "paid"
	mockOrderRepo.On("UpdateOrderStatus", ctx, sessionId, status).Return(nil)
	actual := orderService.UpdateOrderStatus(ctx, sessionId, status)
	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}
