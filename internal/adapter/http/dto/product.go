package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductRequest struct {
	ProductID       primitive.ObjectID `json:"productid"`
	ProductName     string             `json:"productname"`
	Details         string             `json:"details"`
	Stock           int                `json:"stock"`
	Category        string             `json:"category"`
	Location        string             `json:"location"`
	Image           []string           `json:"image"`
	Price           float64            `json:"price"`
	PriceId         string             `json:"priceid"`
	StripeProductId string             `json:"stripeproductid"`
}
type Product struct {
	ProductId     primitive.ObjectID `json:"productid" validate:"required"`
	ProductName   string             `json:"productname" validate:"required"`
	Images        []string           `json:"image" validate:"required"`
	Details       string             `json:"details" validate:"required"`
	Amount        int                `json:"amount" validate:"required"`
	PricePerPiece float64            `json:"priceeach" validate:"required"`
	PriceId       string             `json:"priceid" validate:"required"`
}
