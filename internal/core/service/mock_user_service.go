package service

import (
	"context"
	"golearning/internal/core/domain"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(ctx context.Context, userId string) (*domain.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*domain.User), nil
}
func (m *MockUserService) Login(ctx context.Context, username, password string) (*domain.Token, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(*domain.Token), nil
}
func (m *MockUserService) Register(ctx context.Context, info domain.User) error {
	args := m.Called(ctx, info)
	return args.Error(0)
}
func (m *MockUserService) UpdateUser(ctx context.Context, user domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUserService) ChangePassword(ctx context.Context, userId, oldPassword, newPassword string) error {
	args := m.Called(ctx, userId, oldPassword, newPassword)
	return args.Error(0)
}
func (m *MockUserService) ResetPassword(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), nil
}
func (m *MockUserService) RefreshToken(ctx context.Context, token string) (*domain.Token, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*domain.Token), nil
}
func (m *MockUserService) GetCart(ctx context.Context, userId string) ([]domain.Cart, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]domain.Cart), nil
}
func (m *MockUserService) AddtoCart(ctx context.Context, userProduct domain.Cart, userId string) error {
	args := m.Called(ctx, userProduct, userId)
	return args.Error(0)
}
func (m *MockUserService) IncreaseCartProduct(ctx context.Context, userId string, productId primitive.ObjectID) error {
	args := m.Called(ctx, userId, productId)
	return args.Error(0)
}
func (m *MockUserService) DecreaseCartProduct(ctx context.Context, userId string, productId primitive.ObjectID) error {
	args := m.Called(ctx, userId, productId)
	return args.Error(0)
}
func (m *MockUserService) DeleteItemInCart(ctx context.Context, userId string, productId primitive.ObjectID) error {
	args := m.Called(ctx, userId, productId)
	return args.Error(0)
}
func (m *MockUserService) ClearCart(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}
func (m *MockUserService) GenerateToken(ctx context.Context, username string, id primitive.ObjectID) (*domain.Token, error) {
	args := m.Called(ctx, username, id)
	return args.Get(0).(*domain.Token), nil
}
