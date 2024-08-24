package port

import (
	"context"
	"golearning/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	UpdateStock(ctx context.Context, productId primitive.ObjectID, amount int) error
	AddNewProduct(ctx context.Context, role string, product domain.Product) error
	EditProduct(ctx context.Context, role string, product domain.Product) error
	DeleteProduct(ctx context.Context, role string, product domain.Product) error
	GetAllProduct(context.Context) ([]domain.Product, error)
	GetProductById(ctx context.Context, productId string) (*domain.Product, error)
	CheckAmount(ctx context.Context, productId primitive.ObjectID) (*int, error)
	//ProductCategory(domain.Product) ([]domain.ProductList, error)
}

type ProductRepo interface {
	UpdateStock(ctx context.Context, productId primitive.ObjectID, amount int) error
	CheckAmount(ctx context.Context, productId primitive.ObjectID) (*int, error)
	AddNewProduct(ctx context.Context, product domain.Product) error
	EditProduct(ctx context.Context, product domain.Product) (string, error)
	DeleteProduct(ctx context.Context, productId primitive.ObjectID) error
	GetAllProduct(context.Context) ([]domain.Product, error)
	GetProductById(ctx context.Context, productId primitive.ObjectID) (*domain.Product, error)
}

type ProductCached interface {
	SetProduct(ctx context.Context, json []byte) error
	GetProduct(context.Context) (string, error)
}

//aggregrate acrt objectID to user
