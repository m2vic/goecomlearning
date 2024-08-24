package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProductCached struct {
	Redis *redis.Client
}

func NewProductCache(Redis *redis.Client) *ProductCached {
	return &ProductCached{Redis: Redis}
}

func (r *ProductCached) SetProduct(ctx context.Context, json []byte) error {
	exp := time.Duration(time.Minute * 1)
	err := r.Redis.Set(ctx, "productlist", json, exp)
	if err != nil {
		return fmt.Errorf("error saving to cache")
	}
	return nil
}
func (r *ProductCached) GetProduct(ctx context.Context) (string, error) {
	val, err := r.Redis.Get(ctx, "productlist").Result()
	if err != nil {
		return "", err
	}
	fmt.Println("from cache!")
	return val, nil
}
