package repo

import (
	"context"
	"fmt"
	"golearning/internal/core/domain"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/price"
	"github.com/stripe/stripe-go/v79/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createStripeProduct(item domain.Product) (string, error) {
	env := godotenv.Load()
	if env != nil {
		return "", fmt.Errorf("fail to load env")
	}
	stripeKey := os.Getenv("STRIPEKEY")
	stripe.Key = stripeKey
	ProductParams := &stripe.ProductParams{Name: stripe.String(item.ProductName), Images: stripe.StringSlice(item.Images), Description: &item.Details}
	result, err := product.New(ProductParams)
	if err != nil {
		return "", fmt.Errorf("stripe:%w", err)
	}
	return result.ID, nil
}

func createNewStripePrice(item domain.Product, stripeProductId string) (string, error) {
	env := godotenv.Load()
	if env != nil {
		return "", fmt.Errorf("fail to load env")
	}
	stripeKey := os.Getenv("STRIPEKEY")
	stripe.Key = stripeKey
	PriceParams := &stripe.PriceParams{
		Product:    stripe.String(stripeProductId),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		UnitAmount: stripe.Int64(int64(item.Price * 100)),
	}

	makePrice, err := price.New(PriceParams)
	if err != nil {
		return "", fmt.Errorf("price:%w", err)
	}
	return makePrice.ID, nil

}
func genTokenMongo(user string, id primitive.ObjectID, r *MongoUserRepo) (*domain.Token, error) {
	var token domain.Token
	env := godotenv.Load()
	if env != nil {
		log.Fatalf("err loading: %v", env)
	} else {
		fmt.Println("load complete!")
	}
	accessKey := os.Getenv("PASS")
	refreshKey := os.Getenv("REFRESHTOKEN")
	if accessKey == "" {
		fmt.Println("not setted!")
	}
	exp := time.Now().Add(time.Hour * 1).Unix()
	expRT := time.Now().Add(time.Hour * 1).Unix()
	var role string
	if user == "admin" {
		role = "admin"
	}
	setAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"username": user, "userid": id, "role": role, "exp": exp})
	accessTokenString, err := setAccessToken.SignedString([]byte(accessKey))
	if err != nil {
		log.Fatal(err)
	}
	setRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"username": user, "userid": id, "role": role, "exp": expRT})
	refreshTokenString, err := setRefreshToken.SignedString([]byte(refreshKey))
	if err != nil {
		log.Fatal("from refresh:", err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"token": refreshTokenString}}
	_, err = r.col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	token.AccessToken = accessTokenString
	token.RefreshToken = refreshTokenString
	return &token, nil
}
