package service

import (
	"golearning/internal/core/domain"

	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v79"
)

type MockStripeService struct {
	mock.Mock
}

func (m *MockStripeService) CreateSession(req domain.ProductList) (*stripe.CheckoutSession, error) {
	args := m.Called(req)
	return args.Get(0).(*stripe.CheckoutSession), nil
}
