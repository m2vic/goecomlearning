package service_test

import (
	"context"
	"golearning/internal/adapter/repo"
	"golearning/internal/core/domain"
	"golearning/internal/core/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllProduct(t *testing.T) {
	productRepo := new(repo.MockProductRepo)
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	ctx := context.Background()

	service := service.NewProductService(productRepo, productCached, userRepo)
	expectedProducts := []domain.Product{{ProductName: "KATANA", Price: 12}}
	productRepo.On("GetAllProduct").Return(expectedProducts)
	actual, err := service.GetAllProduct(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, actual)
}

func TestAddNewProduct(t *testing.T) {
	productRepo := new(repo.MockProductRepo)
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	ctx := context.Background()

	productService := service.NewProductService(productRepo, productCached, userRepo)
	var product domain.Product
	role := "admin"
	productRepo.On("AddNewProduct", ctx, product).Return(nil)
	actual := productService.AddNewProduct(ctx, role, product)
	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}

func TestGetProductById(t *testing.T) {
	productRepo := new(repo.MockProductRepo)
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	ctx := context.Background()

	productService := service.NewProductService(productRepo, productCached, userRepo)
	expectedProduct := &domain.Product{ProductName: "Katana", Price: 555}
	productId := "66bfc8a3e3ed16d226c66775"
	hex, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return
	}
	productRepo.On("GetProductById", ctx, hex).Return(expectedProduct)
	actual, err := productService.GetProductById(ctx, productId)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, actual)
}

func TestDeleteProduct(t *testing.T) {
	productRepo := new(repo.MockProductRepo)
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	ctx := context.Background()

	productService := service.NewProductService(productRepo, productCached, userRepo)

	product := domain.Product{ProductID: primitive.NewObjectID()}
	role := "admin"
	productRepo.On("DeleteProduct", ctx, product.ProductID).Return(nil)
	userRepo.On("DeleteItemFromSystem", ctx, product.ProductID).Return(nil)
	actual := productService.DeleteProduct(ctx, role, product)
	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}
func TestEditProduct(t *testing.T) {
	productRepo := new(repo.MockProductRepo)
	userRepo := new(repo.MockUserRepo)
	productCached := new(repo.MockProductCached)
	ctx := context.Background()

	productService := service.NewProductService(productRepo, productCached, userRepo)

	product := domain.Product{ProductID: primitive.NewObjectID()}
	role := "admin"
	productRepo.On("EditProduct", ctx, product).Return(nil)
	userRepo.On("EditItemFromSystem", ctx, product).Return(nil)
	actual := productService.EditProduct(ctx, role, product)
	assert.NoError(t, nil)
	assert.Equal(t, nil, actual)
}
