package service

import (
	"context"
	"golearning/internal/core/domain"
	"golearning/internal/core/port"
)

type OrderService struct {
	OrderRepo port.OrderRepo
}

func NewOrderService(orderRepo port.OrderRepo) *OrderService {
	return &OrderService{OrderRepo: orderRepo}
}

func (s *OrderService) NewOrder(ctx context.Context, order domain.Order) error {

	return s.OrderRepo.NewOrder(ctx, order)
}

func (s *OrderService) GetOrder(ctx context.Context, userId string) ([]domain.Order, error) {

	return s.OrderRepo.GetOrder(ctx, userId)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, sessionId, status string) error {
	return s.OrderRepo.UpdateOrderStatus(ctx, sessionId, status)
}
