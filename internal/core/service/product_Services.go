package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golearning/internal/core/domain"
	"golearning/internal/core/port"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	productRepo  port.ProductRepo
	productCache port.ProductCached
	userRepo     port.UserRepo
}

func NewProductService(
	productRepo port.ProductRepo,
	productCache port.ProductCached,
	userRepo port.UserRepo,
) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		productCache: productCache,
		userRepo:     userRepo,
	}
}

func (r *ProductService) GetAllProduct(ctx context.Context) ([]domain.Product, error) {

	result := []domain.Product{}
	product, err := r.productCache.GetProduct(ctx)
	if err != nil {
		fmt.Println("No Cache!")
	}

	// if cache hit miss
	if product == "" {
		resultFromDb, err := r.productRepo.GetAllProduct(ctx)
		if err != nil {
			return nil, fmt.Errorf("err:%w", err)
		}
		json, err := json.Marshal(resultFromDb)
		if err != nil {
			fmt.Println("fail to marshall")
		}
		_ = r.productCache.SetProduct(ctx, json)
		return resultFromDb, nil
	}
	err = json.Unmarshal([]byte(product), &result)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshall")
	}
	return result, nil
}
func (r *ProductService) GetProductById(ctx context.Context, strProductId string) (*domain.Product, error) {
	productId, err := primitive.ObjectIDFromHex(strProductId)
	if err != nil {
		return nil, err
	}
	return r.productRepo.GetProductById(ctx, productId)
}
func (r *ProductService) AddNewProduct(ctx context.Context, role string, product domain.Product) error {
	if role != "admin" {
		return errors.New("Unauthorized")
	}
	return r.productRepo.AddNewProduct(ctx, product)

}

func (r *ProductService) EditProduct(ctx context.Context, role string, product domain.Product) error {
	if role != "admin" {
		return fmt.Errorf("not Admin")
	}
	newPriceId, err := r.productRepo.EditProduct(ctx, product)
	if err != nil {
		return err
	}
	product.PriceId = newPriceId
	return r.userRepo.EditItemFromSystem(ctx, product)
}

func (r *ProductService) DeleteProduct(ctx context.Context, role string, product domain.Product) error {
	if role != "admin" {
		return fmt.Errorf("not Admin")
	}
	productId := product.ProductID
	err := r.productRepo.DeleteProduct(ctx, productId)
	if err != nil {
		return fmt.Errorf("err from service:%w", err)
	}
	err = r.userRepo.DeleteItemFromSystem(ctx, productId)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductService) CheckAmount(ctx context.Context, productId primitive.ObjectID) (int, error) {
	result, err := r.productRepo.CheckAmount(ctx, productId)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *ProductService) UpdateStock(ctx context.Context, products []domain.StripeProduct) error {
	// we can optimize it []domain.product
	return r.productRepo.UpdateStock(ctx, products)
}
