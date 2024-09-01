package service

import (
	"context"
	"golearning/internal/core/domain"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateToken(ctx context.Context, username string, id primitive.ObjectID) (*domain.Token, error) {
	args := m.Called(ctx, username, id)
	return args.Get(0).(*domain.Token), nil
}
