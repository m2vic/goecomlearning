package repo

import (
	"errors"
	"golearning/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type ProductRepoMock struct {
	mock.Mock
}

func NewMockProductRepo() *ProductRepoMock {
	return &ProductRepoMock{}
}

func (m *ProductRepoMock) GetAllProduct() ([]domain.Product, error) {
	args := m.Called()
	products, ok := args.Get(0).([]domain.Product)
	if !ok {
		return nil, errors.New("unexpected return type")
	}
	return products, nil
}
