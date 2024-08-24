package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cart struct {
	ProductId     primitive.ObjectID `json:"productid" bson:"productid"`
	ProductName   string             `json:"productname" bson:"productname"`
	Images        []string           `json:"image" bson:"image"`
	Details       string             `json:"details" bson:"details"`
	Amount        int                `json:"amount" bson:"amount"`
	PricePerPiece float64            `json:"priceeach" bson:"priceeach"`
	PriceId       string             `json:"priceid" bson:"priceid"`
	//TotalPrice    float64            `json:"totalprice" bson:"totalprice"`
}
