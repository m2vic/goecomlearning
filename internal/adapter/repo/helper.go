package repo

import (
	"fmt"
	"golearning/internal/core/domain"
	"os"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/price"
	"github.com/stripe/stripe-go/v79/product"
)

func createStripeProduct(item domain.Product) (string, error) {
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
