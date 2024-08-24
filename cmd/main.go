package main

import (
	"context"
	"fmt"
	"golearning/internal/adapter/http/handler"
	"golearning/internal/adapter/repo"
	"golearning/internal/core/service"
	"golearning/internal/server"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {

	env := godotenv.Load()
	if env != nil {
		fmt.Println("fail to load env")
	}
	mongoUser := os.Getenv("MONGOUSER")
	mongoPass := os.Getenv("MONGOPASS")
	mongoHost := os.Getenv("MONGOHOST")
	mongoPort := os.Getenv("MONGOPORT")
	mongoDB := os.Getenv("MONGODATABASE")
	mongoAuth := os.Getenv("MONGOAUTH")
	dbName := os.Getenv("DBNAME")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", mongoUser, mongoPass, mongoHost, mongoPort, mongoDB, mongoAuth)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	pong, err := redisClient.Ping(ctx).Result()
	fmt.Println(pong, err)

	stripekey := os.Getenv("STRIPEKEY")
	stripeService := service.NewStripeService(stripekey)
	userRepo := repo.NewMongoRepo(client, dbName, "users")
	orderRepo := repo.NewOrderRepo(client, dbName, "orders")
	orderService := service.NewOrderService(orderRepo)

	productRepo := repo.NewProductRepo(client, "test", "product")

	productCache := repo.NewProductCache(redisClient)
	productService := service.NewProductService(productRepo, productCache, userRepo)
	userService := service.NewUserService(userRepo, productService)
	checkoutService := service.NewCheckoutService(orderService, *stripeService, productService, userService)
	userHandler := handler.NewUserHandler(userService, orderService, checkoutService)
	productHandler := handler.NewProductHandler(*productService)

	server.Start(userHandler, productHandler)
}
