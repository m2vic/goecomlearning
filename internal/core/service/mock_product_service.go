package service

import (
	"context"
	"golearning/internal/core/domain"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) AddNewProduct(ctx context.Context, role string, product domain.Product) error {
	return nil
}
func (m *MockProductService) EditProduct(ctx context.Context, role string, product domain.Product) error {
	return nil
}
func (m *MockProductService) DeleteProduct(ctx context.Context, role string, product domain.Product) error {
	args := m.Called(ctx, role, product)

	return args.Error(0)
}
func (m *MockProductService) GetAllProduct(ctx context.Context) ([]domain.Product, error) {
	args := m.Called(ctx)

	return args.Get(0).([]domain.Product), nil
}
func (m *MockProductService) GetProductById(ctx context.Context, productId string) (*domain.Product, error) {
	args := m.Called(ctx, productId)
	product := args.Get(0).(domain.Product)
	return &product, nil
}

func (m *MockProductService) CheckAmount(ctx context.Context, productId primitive.ObjectID) (*int, error) {
	args := m.Called(ctx, productId)
	amount := args.Int(0)
	return &amount, nil
}

func (m *MockProductService) UpdateStock(ctx context.Context, productId primitive.ObjectID, amount int) error {
	args := m.Called(ctx, productId, amount)
	return args.Error(0)
}
