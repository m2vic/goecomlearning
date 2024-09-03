package port

import (
	"context"
	"golearning/internal/core/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckoutService interface {
	Checkout(ctx context.Context, list domain.ProductList, userId primitive.ObjectID, now time.Time) (string, error)
}
