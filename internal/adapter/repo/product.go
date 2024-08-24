package repo

import (
	"context"
	"errors"
	"fmt"
	"golearning/internal/core/domain"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepo struct {
	client *mongo.Client
	col    *mongo.Collection
}

func NewProductRepo(client *mongo.Client, dbName string, colName string) *ProductRepo {
	if client == nil {
		fmt.Println("missing mongodb client")
		return nil
	}

	col := client.Database(dbName).Collection(colName)
	return &ProductRepo{client: client, col: col}
}

func (r *ProductRepo) GetAllProduct(ctx context.Context) ([]domain.Product, error) {
	cur, err := r.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to find products:%w", err)
	}
	defer cur.Close(ctx)

	var products []domain.Product

	for cur.Next(ctx) {
		var product domain.Product
		err := cur.Decode(&product)
		if err != nil {
			return nil, fmt.Errorf("failed to decode document:%w", err)
		}
		products = append(products, product)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error:%w", err)
	}

	return products, nil
}
func (r *ProductRepo) GetProductById(ctx context.Context, productId primitive.ObjectID) (*domain.Product, error) {
	var result domain.Product
	filter := bson.M{"_id": productId}
	err := r.col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("error from repo:%w", err)
	}
	return &result, nil
}
func (r *ProductRepo) AddNewProduct(ctx context.Context, product domain.Product) error {
	stripeProductId, err := createStripeProduct(product)
	if err != nil {
		return fmt.Errorf("stripe:%w", err)
	}
	priceId, err := createNewStripePrice(product, stripeProductId)
	if err != nil {
		return fmt.Errorf("stripe:%w", err)
	}
	productId := primitive.NewObjectID()
	product.ProductID = productId
	product.PriceId = priceId
	product.StripeProductId = stripeProductId
	_, err = r.col.InsertOne(ctx, product)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
func (r *ProductRepo) EditProduct(ctx context.Context, product domain.Product) (string, error) {
	priceId, err := createNewStripePrice(product, product.StripeProductId)
	if err != nil {
		fmt.Println("fail to create a new stripe price")
		return "", err
	}
	filter := bson.M{"_id": product.ProductID}

	update := bson.M{"$set": bson.M{"productname": product.ProductName,
		"category": product.Category,
		"details":  product.Details, "stock": product.Stock, "price": product.Price,
		"priceid": priceId, "location": product.Location}}

	_, err = r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("fail to update to db")
		return "", err
	}
	return priceId, nil
}
func (r *ProductRepo) DeleteProduct(ctx context.Context, productId primitive.ObjectID) error {

	filter := bson.M{"_id": productId}
	_, err := r.col.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("fail to delete at repo layer:%w", err)
	}
	return nil
}
func (r *ProductRepo) CheckAmount(ctx context.Context, productId primitive.ObjectID) (*int, error) {
	var product domain.Product
	filter := bson.M{"_id": productId}
	err := r.col.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product.Stock, nil

}

func (r *ProductRepo) UpdateStock(ctx context.Context, productId primitive.ObjectID, amount int) error {
	filter := bson.M{"_id": productId}
	update := bson.M{"$inc": bson.M{"stock": -amount}}
	_, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Error:%w", err)
	}
	return nil
}
