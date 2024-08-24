package domain

import "time"

type Order struct {
	UserId         string           `json:"userid" bson:"userid"`
	OrderID        string           `json:"orderid" bson:"orderid"` // store checkout session id from stripe
	Ordered_At     time.Time        `json:"ordered_at" bson:"ordered_at"`
	TotalPrice     int              `json:"totalprice" bson:"totalprice"`
	Discount       int              `json:"discount" bson:"discount"`
	Payment_Method string           `json:"payment_method" bson:"payment_method"`
	Status         string           `json:"status" bson:"status"`
	LineItems      []ProductDetails `json:"line_items" bson:"line_items"`
}

type ProductDetails struct {
	ProductName string   `json:"productname"`
	Description string   `json:"description"`
	Quantity    int      `json:"quantity"`
	UnitPrice   int      `json:"unitprice"`
	Images      []string `json:"images"`
}
