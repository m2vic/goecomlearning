package port

import (
	"golearning/internal/core/domain"

	"github.com/stripe/stripe-go/v79"
)

type StripeService interface {
	CreateSession(req domain.ProductList) (*stripe.CheckoutSession, error)
}
