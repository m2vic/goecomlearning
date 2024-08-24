package port

import (
	"context"
	"golearning/internal/core/domain"
)

type OrderService interface {
	NewOrder(ctx context.Context, order domain.Order) error
	GetOrder(ctx context.Context, userId string) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, sessionId, status string) error
}

type OrderRepo interface {
	NewOrder(ctx context.Context, order domain.Order) error
	GetOrder(ctx context.Context, userId string) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, sessionId, status string) error
}

//what if one user have many order
// delete order after is xx years old
