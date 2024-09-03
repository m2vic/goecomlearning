package repo

import (
	"context"
	"golearning/internal/core/domain"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUser(ctx context.Context, userId primitive.ObjectID) (*domain.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*domain.User), nil
}

func (m *MockUserRepo) Register(ctx context.Context, username, email string, password []byte) error {
	args := m.Called(ctx, username, email, password)
	return args.Error(0)
}
func (m *MockUserRepo) UpdateUser(ctx context.Context, info domain.User) error {
	args := m.Called(ctx, info)
	return args.Error(0)
}
func (m *MockUserRepo) CheckEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), nil
}
func (m *MockUserRepo) CheckUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), nil
}

func (m *MockUserRepo) ResetPassword(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), nil
}
func (m *MockUserRepo) ChangePassword(ctx context.Context, userId primitive.ObjectID, newPassword string) error {
	args := m.Called(ctx, userId, newPassword)
	return args.Error(0)
}
func (m *MockUserRepo) CheckRefresh(ctx context.Context, token string) (*domain.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*domain.User), nil
}
func (m *MockUserRepo) GetCart(ctx context.Context, userId primitive.ObjectID) ([]domain.Cart, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]domain.Cart), nil
}
func (m *MockUserRepo) AddtoCart(ctx context.Context, userProduct domain.Cart, userId primitive.ObjectID) error {
	args := m.Called(ctx, userProduct, userId)
	return args.Error(0)
}
func (m *MockUserRepo) IncreaseCartProduct(ctx context.Context, userId primitive.ObjectID, productId primitive.ObjectID) error {
	args := m.Called(ctx, userId, productId)
	return args.Error(0)
}
func (m *MockUserRepo) DecreaseCartProduct(ctx context.Context, userId primitive.ObjectID, productId primitive.ObjectID) error {
	args := m.Called(ctx, userId, productId)
	return args.Error(0)
}
func (m *MockUserRepo) DeleteItemInCart(ctx context.Context, userId primitive.ObjectID, ProductId primitive.ObjectID) error {
	args := m.Called(ctx, userId, ProductId)
	return args.Error(0)
}
func (m *MockUserRepo) DeleteItemFromSystem(ctx context.Context, productId primitive.ObjectID) error {
	args := m.Called(ctx, productId)
	return args.Error(0)
}
func (m *MockUserRepo) EditItemFromSystem(ctx context.Context, product domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}
func (m *MockUserRepo) ClearCart(ctx context.Context, userId primitive.ObjectID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}
func (m *MockUserRepo) CheckPassword(ctx context.Context, oldPass string) error {
	args := m.Called(ctx, oldPass)
	return args.Error(0)
}
