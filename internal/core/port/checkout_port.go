package port

import (
	"context"
	"golearning/internal/core/domain"
	"time"
)

type CheckoutService interface {
	Checkout(ctx context.Context, list domain.ProductList, userId string, now time.Time) (string, error)
}
