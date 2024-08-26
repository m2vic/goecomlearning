package repo

import (
	"context"
	"golearning/internal/core/domain"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) AddNewProduct(ctx context.Context, product domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}
func (m *MockProductRepo) EditProduct(ctx context.Context, product domain.Product) (string, error) {
	return "", nil
}
func (m *MockProductRepo) DeleteProduct(ctx context.Context, productId primitive.ObjectID) error {
	args := m.Called(ctx, productId)
	return args.Error(0)
}
func (m *MockProductRepo) GetAllProduct(context.Context) ([]domain.Product, error) {
	args := m.Called()
	return args.Get(0).([]domain.Product), nil
}
func (m *MockProductRepo) GetProductById(ctx context.Context, productId primitive.ObjectID) (*domain.Product, error) {
	args := m.Called(ctx, productId)
	return args.Get(0).(*domain.Product), nil
}
func (m *MockProductRepo) CheckAmount(ctx context.Context, productId primitive.ObjectID) (*int, error) {
	args := m.Called(ctx, productId)
	amount := args.Int(0)
	return &amount, nil
}
func (m *MockProductRepo) UpdateStock(ctx context.Context, product []domain.StripeProduct) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}
