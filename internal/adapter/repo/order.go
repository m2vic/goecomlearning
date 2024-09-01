package repo

import (
	"context"
	"errors"
	"fmt"
	"golearning/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepo struct {
	client *mongo.Client
	col    *mongo.Collection
}

func NewOrderRepo(client *mongo.Client, dbName string, colName string) *OrderRepo {
	if client == nil {
		fmt.Println("missing client")
		return nil
	}
	col := client.Database(dbName).Collection(colName)
	return &OrderRepo{client: client, col: col}
}

func (r *OrderRepo) NewOrder(ctx context.Context, order domain.Order) error {
	_, err := r.col.InsertOne(ctx, order)
	if err != nil {
		return fmt.Errorf("from db:%w", err)
	}
	return nil
}

func (r *OrderRepo) GetOrder(ctx context.Context, userId string) ([]domain.Order, error) {
	orders := []domain.Order{}
	filter := bson.M{"userid": userId}
	//have  handle many orderctx context.Context,
	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("from repo:%w", err)
		} else {
			return nil, fmt.Errorf("from repo:%w", err)
		}
	}
	for cur.Next(context.Background()) {
		order := domain.Order{}
		err := cur.Decode(&order)
		if err != nil {
			return nil, errors.New("fail to map orders")
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepo) UpdateOrderStatus(ctx context.Context, sessionId, status string) error {
	filter := bson.M{"orderid": sessionId}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("fail to update.")
	}
	return nil
}
