package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ProductID       primitive.ObjectID `bson:"_id" json:"productid"`
	ProductName     string             `bson:"productname" json:"productname"`
	Details         string             `bson:"details" json:"details"`
	Stock           int                `bson:"stock" json:"stock"`
	Category        string             `bson:"category" json:"category"`
	Location        string             `bson:"location" json:"location"`
	Images          []string           `bson:"image" json:"image"`
	Price           float64            `bson:"price" json:"price"`
	PriceId         string             `bson:"priceid" json:"priceid"`
	StripeProductId string             `bson:"stripeproductid" json:"stripeproductid"`
}

type UsersProduct struct {
	ProductID     primitive.ObjectID `bson:"productid" json:"productid"`
	ProductName   string             `bson:"productname" json:"productname"`
	Amount        int                `bson:"amount" json:"amount"`
	PricePerPiece float64            `bson:"priceeach" json:"priceeach"`
	PriceId       string             `bson:"priceid" json:"priceid"`
	//TotalPrice    float64            `bson:"totalprice" json:"totalprice"`
}

type ProductList struct {
	ProductList []StripeProduct `json:"productlist"`
}

type StripeProduct struct {
	ProductId     primitive.ObjectID `json:"productid"`
	ProductName   string             `json:"productname" validate:"required"`
	Images        []string           `json:"image"`
	Details       string             `json:"details"`
	Amount        int                `json:"amount" validate:"required"`
	PricePerPiece float64            `json:"priceeach" validate:"required"`
	PriceId       string             `json:"priceid" validate:"required"`
}

func PriceCal(amount int, price float64) float64 {
	return price * float64(amount)
}
